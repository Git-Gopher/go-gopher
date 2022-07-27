package main

import (
	"fmt"
	"log"

	"github.com/Git-Gopher/go-gopher/assess"
	"github.com/Git-Gopher/go-gopher/model"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/Git-Gopher/go-gopher/version"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/urfave/cli/v2"
)

func singleCommand(cCtx *cli.Context) error {
	fmt.Printf("BuildVersion: %v\n", version.BuildVersion())
	utils.Environment(".env")
	// Handle flags.
	githubURL := cCtx.String("url")

	// Clone repository into memory.
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: githubURL,
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

	candidates := assess.RunMarker(
		assess.MarkerCtx{
			Model:        enrichedModel,
			Contribution: assess.NewContribution(*enrichedModel),
			Author:       authors,
		},
		assess.BasicGradingAlgorithm,
		// Markers
		assess.D(assess.Atomicity),
		assess.D(assess.CommitMessage),
		assess.D(assess.RegularBranchNames),
		assess.D(assess.FeatureBranching),
		assess.D(assess.PullRequestReview),
	)

	for _, candidate := range candidates {
		log.Printf("#### @%s ####\n", candidate.Username)
	}

	if err := IndividualReports(candidates); err != nil {
		return fmt.Errorf("failed to generate individual reports: %w", err)
	}

	return nil
}
