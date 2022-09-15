package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Git-Gopher/go-gopher/cache"
	"github.com/Git-Gopher/go-gopher/discord"
	"github.com/Git-Gopher/go-gopher/markup"
	"github.com/Git-Gopher/go-gopher/model"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/Git-Gopher/go-gopher/version"
	"github.com/Git-Gopher/go-gopher/violation"
	"github.com/Git-Gopher/go-gopher/workflow"
	"github.com/go-git/go-git/v5"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var errOwnerMismatch = errors.New("owner mismatch")

func actionCommand(cCtx *cli.Context) error {
	log.Printf("BuildVersion: %s", version.BuildVersion())
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
	caches, err := cache.Read()
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to read caches: %w", err)
		}

		// Write a cache for current so that next run can use it.
		if err = cache.Write([]*cache.Cache{current}); err != nil {
			return fmt.Errorf("failed to write cache: %w", err)
		}
	}

	cfg := readConfig(cCtx)
	ghwf := workflow.GithubFlowWorkflow(cfg)
	violated, count, total, violations, err := ghwf.Analyze(enrichedModel, authors, current, caches)
	if err != nil {
		log.Fatalf("Failed to analyze: %v\n", err)
	}

	if config.LoginWhiteList != "" {
		whitelist := strings.Split(config.LoginWhiteList, ",")
		violations = violation.FilterByLogin(violations, whitelist)
	}

	disallowList := []string{"github-classroom[bot]"}
	filteredViolations := violation.FilterByLogin(violations, disallowList)
	workflow.PrintSummary(authors, violated, count, total, filteredViolations)

	summary := workflow.MarkdownSummary(authors, filteredViolations)
	markup.Outputs("pr_summary", summary)

	fn, err := ghwf.WriteLog(*enrichedModel, cfg)
	if err != nil {
		log.Printf("Could not write json log: %v", err)

		return nil
	}

	if err = discord.SendLog(fn); err != nil {
		log.Printf("Could not write json log to discord: %v", err)

		return nil
	}

	return nil
}
