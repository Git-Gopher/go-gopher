package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Git-Gopher/go-gopher/cache"
	"github.com/Git-Gopher/go-gopher/markup"
	"github.com/Git-Gopher/go-gopher/model"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/Git-Gopher/go-gopher/violation"
	"github.com/Git-Gopher/go-gopher/workflow"
	"github.com/go-git/go-git/v5"
	"github.com/urfave/cli/v2"
)

var errOwnerMismatch = errors.New("owner mismatch")

func actionCommand(cCtx *cli.Context) error {
	// Load the environment variables from GitHub Actions.
	config, err := loadEnv(cCtx.Context)
	if err != nil {
		return fmt.Errorf("failed to load env: %w", err)
	}

	// Open the repository.
	repo, err := git.PlainOpen(config.GithubWorkspace)
	if err != nil {
		return fmt.Errorf("failed to open repo: %w", err)
	}

	// GithubURL fallback.
	githubURL, err := utils.Url(repo)
	if err != nil {
		return fmt.Errorf("failed to get url: %w", err)
	}

	// Get the repositoryName.
	repoOwner, repoName, err := utils.OwnerNameFromUrl(githubURL)
	if err != nil {
		return fmt.Errorf("failed to get owner and repo name: %w", err)
	}
	if config.GithubRepositoryOwner != repoOwner {
		return fmt.Errorf("%w: %s != %s", errOwnerMismatch, repoOwner, config.GithubRepositoryOwner)
	}

	// Create enrichedModel.
	enrichedModel, err := model.FetchEnrichedModel(repo, repoOwner, repoName)
	if err != nil {
		return fmt.Errorf("failed to create enriched model: %w", err)
	}

	// Create cache.
	current := cache.NewCache(enrichedModel)

	// Populate authors from enrichedModel.
	authors := enriched.PopulateAuthors(enrichedModel)

	// Read cache.
	caches, err := cache.ReadCaches()
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to read caches: %w", err)
		}

		// Write a cache for current so that next run can use it.
		if err = cache.WriteCaches([]*cache.Cache{current}); err != nil {
			return fmt.Errorf("failed to write cache: %w", err)
		}
	}

	cfg := readConfig(cCtx)
	ghwf := workflow.GithubFlowWorkflow(cfg)
	violated, count, total, violations, err := ghwf.Analyze(enrichedModel, authors, current, caches)
	if err != nil {
		log.Fatalf("Failed to analyze: %v\n", err)
	}

	workflowSummary(authors, violated, count, total, violations)

	if cCtx.Bool("csv") {
		err = ghwf.Csv(workflow.DefaultCsvPath, enrichedModel.Name, enrichedModel.URL)
		if err != nil {
			log.Fatalf("Could not create csv summary: %v", err)
		}
	}

	return nil
}

// Print violation summary to IO, Split by severity with author association.
func workflowSummary(authors utils.Authors, v, c, t int, vs []violation.Violation) {
	var violations, suggestions []violation.Violation
	for _, v := range vs {
		switch v.Severity() {
		case violation.Violated:
			violations = append(violations, v)
		case violation.Suggestion:
			suggestions = append(suggestions, v)
		}
	}

	var vsd string
	for _, v := range violations {
		vsd += v.Display(authors)
	}
	markup.Group("Violations", vsd)

	var ssd string
	for _, v := range suggestions {
		ssd += v.Display(authors)
	}
	markup.Group("Suggestions", ssd)

	var asd string
	counts := make(map[string]int)
	for _, v := range vs {
		email := v.Email()
		login, err := authors.Find(email)
		if err != nil {
			continue
		}
		counts[*login]++
	}

	for login, count := range counts {
		asd += fmt.Sprintf("%s: %d\n", login, count)
	}

	asd += fmt.Sprintf("violated: %d\n", v)
	asd += fmt.Sprintf("count: %d\n", c)
	asd += fmt.Sprintf("total: %d\n", t)
	markup.Group("Summary", asd)
}