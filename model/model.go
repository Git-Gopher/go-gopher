package model

import (
	"fmt"
	"time"

	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/model/remote"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	log "github.com/sirupsen/logrus"
)

//nolint:gocognit
func FetchEnrichedModel(repo *git.Repository, repoOwner, repoName string) (*enriched.EnrichedModel, error) {
	log.Infof("Beginning to scrape repository %s", repoName)
	// scraping remote GitHub repository.
	start := time.Now()

	githubModel, err := remote.ScrapeRemoteModel(repoOwner, repoName)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape remote model: %w", err)
	}

	elapsed := time.Since(start)
	log.Infof("Scraped remote GitHub repository in %s", elapsed)

	// here we need to define which is the dev branch and which is release branch
	// most of the time we assume that the dev branch is the default branch

	// dev branch:
	// we can determine this by the most number of PR merged into the branch is the dev branch
	devBranch := findDevBranchByPR(githubModel.PullRequests)

	// release branch:
	// we can determine this by the most number of tags is the release branch

	// loading local Git repository.
	start = time.Now()

	gitModel, err := local.NewGitModel(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to create local model: %w", err)
	}
	elapsed = time.Since(start)
	log.Infof("Loaded local Git repository in %s", elapsed)

	enrichedModel := enriched.NewEnrichedModel(*gitModel, *githubModel)

	// if the main branch isn't the dev branch, we need to find the dev branch
	if devBranch != nil && gitModel.MainGraph.BranchName != *devBranch {
		enrichedModel.ReleaseGraph = gitModel.MainGraph // the main branch is the release branch
		var headCommit *local.Hash
		for _, branch := range gitModel.Branches {
			// remote branch name has origin/ as prefix
			if branch.Name == *devBranch {
				h := branch.Head.Hash
				headCommit = &h

				break
			}
		}

		if headCommit != nil {
			refCommit, err := repo.CommitObject(plumbing.NewHash(headCommit.HexString()))

			log.Infof("devBranch: %s", *devBranch)
			log.Infof("refCommit: %s", headCommit.HexString())

			if err != nil {
				return nil, fmt.Errorf("failed to get commit object: %w", err)
			}

			enrichedModel.MainGraph = local.FetchBranchGraph(refCommit)
			enrichedModel.MainGraph.BranchName = *devBranch
		} else {
			enrichedModel.ReleaseGraph = nil
		}
	}

	// if the default branch is not release branch, try find release branch
	if enrichedModel.ReleaseGraph == nil { //nolint:nestif
		tags, err := repo.Tags()
		if err != nil {
			return nil, fmt.Errorf("failed to get tags: %w", err)
		}

		revHashMap := make(map[string]struct{})

		err = tags.ForEach(func(t *plumbing.Reference) error {
			// This technique should work for both lightweight and annotated tags.
			revHash, err2 := repo.ResolveRevision(plumbing.Revision(t.Name()))
			if err != nil {
				return fmt.Errorf("failed to resolve revision: %w", err2)
			}

			if revHash == nil {
				return nil
			}

			revHashMap[revHash.String()] = struct{}{}

			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to iterate tags: %w", err)
		}

		for _, branch := range enrichedModel.Branches {
			if _, ok := revHashMap[branch.Head.Hash.HexString()]; ok {
				if branch.Name == gitModel.MainGraph.BranchName {
					continue
				}

				refCommit, err := repo.CommitObject(plumbing.NewHash(branch.Head.Hash.HexString()))
				if err != nil {
					return nil, fmt.Errorf("failed to get commit object: %w", err)
				}

				enrichedModel.ReleaseGraph = local.FetchBranchGraph(refCommit)

				break
			}
		}
	}

	log.Infof("Main/Dev Branch: %s", enrichedModel.MainGraph.BranchName)
	if enrichedModel.ReleaseGraph != nil {
		log.Infof("Release Branch: %s", enrichedModel.ReleaseGraph.BranchName)
	}

	return enrichedModel, nil
}

// findDevBranchByPR finds the dev branch by the most number of PR merged into the branch.
func findDevBranchByPR(prs []*remote.PullRequest) *string {
	if len(prs) == 0 {
		return nil
	}

	count := map[string]int{}
	for _, pr := range prs {
		if !pr.Merged {
			continue
		}

		if _, ok := count[pr.BaseRefName]; !ok {
			count[pr.BaseRefName] = 0
		}
		count[pr.BaseRefName]++
	}

	if len(count) == 0 {
		return nil
	}

	max := 0
	var devBranch string
	for branch, c := range count {
		if c > max {
			max = c
			devBranch = branch
		}
	}

	return &devBranch
}
