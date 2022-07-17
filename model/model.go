package model

import (
	"fmt"
	"log"
	"time"

	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/model/remote"
	"github.com/go-git/go-git/v5"
)

func FetchEnrichedModel(repo *git.Repository, repoOwner, repoName string) (*enriched.EnrichedModel, error) {
	// scraping remote GitHub repository.
	start := time.Now()

	githubModel, err := remote.ScrapeRemoteModel(repoOwner, repoName)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape remote model: %w", err)
	}

	elapsed := time.Since(start)
	log.Printf("Scraped remote GitHub repository in %s", elapsed)

	// loading local Git repository.
	start = time.Now()

	gitModel, err := local.NewGitModel(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to create local model: %w", err)
	}
	elapsed = time.Since(start)
	log.Printf("Loaded local Git repository in %s", elapsed)

	enrichedModel := enriched.NewEnrichedModel(*gitModel, *githubModel)

	return enrichedModel, nil
}
