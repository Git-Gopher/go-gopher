package main

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/Git-Gopher/go-gopher/assess"
	"github.com/Git-Gopher/go-gopher/assess/markers/analysis"
	"github.com/Git-Gopher/go-gopher/model"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/Git-Gopher/go-gopher/version"
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

func singleUrlCommand(cCtx *cli.Context) error {
	log.Printf("BuildVersion: %v\n", version.BuildVersion())
	githubURL := cCtx.Args().Get(0)
	if githubURL == "" {
		return errGitHubURL
	}

	var auth *githttp.BasicAuth
	if cCtx.String("token") != "" {
		auth = &githttp.BasicAuth{
			Username: "non-empty",
			Password: cCtx.String("token"),
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

	if err := IndividualReports(candidates); err != nil {
		return fmt.Errorf("failed to generate individual reports: %w", err)
	}

	return nil
}

func singleLocalCommand(cCtx *cli.Context) error {
	log.Printf("BuildVersion: %v\n", version.BuildVersion())

	directory := cCtx.Args().Get(0)
	if directory == "" {
		return errLocalDir
	}

	return runLocalRepository(directory)
}

func runLocalRepository(directory string) error {
	// Open repository locally.
	repo, err := git.PlainOpen(directory)
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	githubURL, err := utils.Url(repo)
	if err != nil {
		return fmt.Errorf("failed to get url: %w", err)
	}

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

	if err := IndividualReports(candidates); err != nil {
		return fmt.Errorf("failed to generate individual reports: %w", err)
	}

	return nil
}

func folderLocalCommand(cCtx *cli.Context) error {
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

	// Load repos into channel.
	go func() {
		for _, repo := range repos {
			wg.Add(1)
			repoChan <- repo
		}
	}()

	for i := 0; i < runtime.NumCPU()-1; i++ {
		go func() {
			select {
			case repo := <-repoChan:
				if err := runLocalRepository(repo); err != nil {
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
