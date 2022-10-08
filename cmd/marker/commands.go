package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/Git-Gopher/go-gopher/assess"
	"github.com/Git-Gopher/go-gopher/assess/markers/analysis"
	"github.com/Git-Gopher/go-gopher/assess/options"
	"github.com/Git-Gopher/go-gopher/model"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	log "github.com/sirupsen/logrus"
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

	if err = c.runMarker(repo, githubURL, flags.LookupPath); err != nil {
		return err
	}

	return nil
}

func (c *Cmds) SingleLocalCommand(cCtx *cli.Context, flags *Flags) error {
	directory := cCtx.Args().Get(0)
	if directory == "" {
		return errLocalDir
	}

	return c.runLocalRepository(directory, flags.LookupPath)
}

func (c *Cmds) runLocalRepository(directory string, lookupPath string) error {
	// Open repository locally.
	repo, err := git.PlainOpen(directory)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	githubURL, err := utils.Url(repo)
	if err != nil {
		return fmt.Errorf("failed to get url: %w", err)
	}

	if err = c.runMarker(repo, githubURL, lookupPath); err != nil {
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
				if err := c.runLocalRepository(repo, flags.LookupPath); err != nil {
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

func (c *Cmds) runMarker(repo *git.Repository, githubURL string, lookupPath string) error {
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

	// Fetch lookup.
	upis, fullnames := fetchLookup(lookupPath)
	// Read marker configs
	o := LoadOptions(log.StandardLogger())
	analyzers := assess.LoadAnalyzer(o)

	cutoff, err := time.Parse("2006-01-02 15:04:05 -0700 MST", o.CutoffDate)
	if err != nil {
		return fmt.Errorf("failed to parse cutoff date %w", err)
	}

	candidates := assess.RunMarker(
		analysis.MarkerCtx{
			Model:        enrichedModel,
			Contribution: analysis.NewContribution(*enrichedModel),
			Author:       authors,
			CutoffDate:   cutoff,
		},
		analyzers,
	)

	for _, candidate := range candidates {
		log.Printf("#### @%s ####\n", candidate.Username)
	}

	if err := IndividualReports(o, repoName, candidates, upis, fullnames); err != nil {
		return fmt.Errorf("failed to generate individual reports: %w", err)
	}

	return nil
}

// Fetch login => (fullname, upi) lookup.
// nolint
func fetchLookup(lookupPath string) (upis map[string]string, fullnames map[string]string) {
	f, err := os.Open(filepath.Clean(lookupPath))
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	// read csv values using csv.Reader
	reader := csv.NewReader(f)
	data, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	upis = make(map[string]string)
	fullnames = make(map[string]string)

	for _, v := range data {
		// unsafe, too bad!
		login := v[2]

		upi := v[3]
		fullname := v[4] + v[5]

		upis[login] = upi
		fullnames[login] = fullname
	}

	return upis, fullnames
}
