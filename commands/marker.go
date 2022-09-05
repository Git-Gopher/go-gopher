package commands

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/Git-Gopher/go-gopher/assess"
	"github.com/Git-Gopher/go-gopher/assess/markers/analysis"
	"github.com/Git-Gopher/go-gopher/assess/options"
	"github.com/Git-Gopher/go-gopher/markup"
	"github.com/Git-Gopher/go-gopher/model"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/gomarkdown/markdown"
	"github.com/urfave/cli/v2"
)

var (
	errGitHubURL = fmt.Errorf("missing GitHub URL")
	errLocalDir  = fmt.Errorf("missing Local Directory")
)

var _ Commands = &Cmds{}

type Commands interface {
	SingleUrlCommand(cCtx *cli.Context, flags *Flags) error
	SingleLocalCommand(cCtx *cli.Context, flags *Flags) error
	FolderLocalCommand(cCtx *cli.Context, flags *Flags) error
	GenerateConfigCommand(cCtx *cli.Context, flags *Flags) error
}

type Cmds struct{}

func (c *Cmds) SingleUrlCommand(cCtx *cli.Context, flags *Flags) error {
	githubURL := cCtx.Args().Get(0)
	if githubURL == "" {
		return errGitHubURL
	}

	var auth *githttp.BasicAuth
	if flags.GithubToken != "" {
		auth = &githttp.BasicAuth{
			Username: "non-empty",
			Password: flags.GithubToken,
		}
	}

	// Clone repository into memory.
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:  githubURL,
		Auth: auth,
	})
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	if err = c.runMarker(repo, githubURL); err != nil {
		return err
	}

	return nil
}

func (c *Cmds) SingleLocalCommand(cCtx *cli.Context, flags *Flags) error {
	directory := cCtx.Args().Get(0)
	if directory == "" {
		return errLocalDir
	}

	return c.runLocalRepository(directory)
}

func (c *Cmds) runLocalRepository(directory string) error {
	// Open repository locally.
	repo, err := git.PlainOpen(directory)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	githubURL, err := utils.Url(repo)
	if err != nil {
		return fmt.Errorf("failed to get url: %w", err)
	}

	if err = c.runMarker(repo, githubURL); err != nil {
		return err
	}

	return nil
}

func (c *Cmds) FolderLocalCommand(cCtx *cli.Context, flags *Flags) error {
	directory := cCtx.Args().Get(0)
	if directory == "" {
		return errLocalDir
	}

	repos := make([]string, 0)

	err := filepath.Walk(directory, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() && info.Name() == ".git" {
			repos = append(repos, filepath.Dir(path))

			return filepath.SkipDir
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	repoChan := make(chan string, runtime.NumCPU())

	wg.Add(1)

	// Load repos into channel.
	go func() {
		for _, repo := range repos {
			wg.Add(1)
			repoChan <- repo
		}
		wg.Done()
	}()

	for i := 0; i < runtime.NumCPU()-1; i++ {
		go func() {
			select {
			case repo := <-repoChan:
				if err := c.runLocalRepository(repo); err != nil {
					log.Errorf("failed to run local repository: %v", err)
				}
				wg.Done()
			case <-ctx.Done():
				return
			}
		}()
	}

	wg.Wait()
	cancel()

	log.Printf("# Done %s #\n", directory)

	return nil
}

func (c *Cmds) GenerateConfigCommand(cCtx *cli.Context, flags *Flags) error {
	r := options.NewFileReader(log.StandardLogger(), nil)
	if err := r.GenerateDefault(flags.OptionsDir); err != nil {
		if !errors.Is(err, utils.ErrSkipped) {
			return fmt.Errorf("failed to generate default options: %w", err)
		}
	}

	if err := utils.GenerateEnv(flags.EnvDir); err != nil {
		if !errors.Is(err, utils.ErrSkipped) {
			return fmt.Errorf("failed to generate env: %w", err)
		}
	}

	return nil
}

func (c *Cmds) runMarker(repo *git.Repository, githubURL string) error {
	// Get the repositoryName.
	repoOwner, repoName, err := utils.OwnerNameFromUrl(githubURL)
	if err != nil {
		return fmt.Errorf("failed to get owner and repo name: %w", err)
	}

	// Create enrichedModel.
	enrichedModel, err := model.FetchEnrichedModel(repo, repoOwner, repoName)
	if err != nil {
		return fmt.Errorf("failed to create enriched model: %w", err)
	}

	// Populate authors from enrichedModel.
	authors := enriched.PopulateAuthors(enrichedModel)

	// Read marker configs
	o := LoadOptions(log.StandardLogger())
	analyzers := assess.LoadAnalyzer(o)

	candidates := assess.RunMarker(
		analysis.MarkerCtx{
			Model:        enrichedModel,
			Contribution: analysis.NewContribution(*enrichedModel),
			Author:       authors,
		},
		analyzers,
	)

	for _, candidate := range candidates {
		log.Printf("#### @%s ####\n", candidate.Username)
	}

	if err := IndividualReports(o, repoName, candidates); err != nil {
		return fmt.Errorf("failed to generate individual reports: %w", err)
	}

	return nil
}

func IndividualReports(options *options.Options, repoName string, candidates []assess.Candidate) error {
	if len(candidates) == 0 {
		return fmt.Errorf("no candidates") //nolint: goerr113
	}

	markers := make([]string, len(candidates[0].Grades))
	for i, grade := range candidates[0].Grades {
		markers[i] = grade.Name
	}

	for _, candidate := range candidates {
		rows := make([][]string, len(candidate.Grades))
		for i, grade := range candidate.Grades {
			rows[i] = []string{
				grade.Name,
				fmt.Sprintf("%d", grade.Violation),
				fmt.Sprintf("%d", grade.Contribution),
				fmt.Sprintf("%d/3", grade.Grade),
			}
		}

		header := []string{"Marker", "Violation", "Contribution", "Grade"}

		md := markup.CreateMarkdown(fillTemplate(options.HeaderTemplate, candidate.Username, repoName)).
			Header("Marked by git-gopher", 2).
			Table(header, rows)

		for _, grade := range candidate.Grades {
			subcategory := fmt.Sprintf("Marker details %s", grade.Name)
			md.Header(subcategory, 2).Paragraph(grade.Details)
		}

		output := markdown.ToHTML([]byte(md.Render()), nil, nil)

		filename := fillTemplate(options.FilenameTemplate, candidate.Username, repoName) + ".html"
		if len(options.OutputDir) != 0 {
			if _, err := os.Stat(options.OutputDir); errors.Is(err, os.ErrNotExist) {
				if err2 := os.MkdirAll(options.OutputDir, os.ModePerm); err2 != nil {
					return fmt.Errorf("can't create options dir: %w", err)
				}
			}
			filename = filepath.Join(options.OutputDir, filename)
		}

		if err := writeFile(filename, output); err != nil {
			return fmt.Errorf("failed to write file for %s: %w", filename, err)
		}
	}

	return nil
}

func fillTemplate(template string, username string, repository string) string {
	template = strings.ReplaceAll(template, "{{.Username}}", username)
	template = strings.ReplaceAll(template, "{{ .Username }}", username)

	template = strings.ReplaceAll(template, "{{.Repository}}", repository)
	template = strings.ReplaceAll(template, "{{ .Repository }}", repository)

	return template
}

func writeFile(path string, data []byte) (err error) {
	f, err := os.Create(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer func() {
		err = f.Close()
	}()

	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", path, err)
	}

	return nil
}
