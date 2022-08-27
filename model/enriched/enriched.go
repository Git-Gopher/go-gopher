package enriched

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/model/remote"
	"github.com/Git-Gopher/go-gopher/utils"
	log "github.com/sirupsen/logrus"
)

var (
	ErrPullRequestNumber = errors.New("could not fetch pull request number from env (PR_NUMBER)")
	ErrFindPullRequest   = errors.New("could not find pull request from scraped repo given pull request number")
)

type EnrichedModel struct {
	// local.GitModel
	Owner           string
	Name            string
	URL             string
	Commits         []local.Commit
	Branches        []local.Branch
	MainGraph       *local.BranchGraph    `json:"-"` // Graph representation of commits in the main branch
	BranchMatrix    []*local.BranchMatrix `json:"-"` // Matrix representation by comparing branches
	LocalCommitters []local.Committer

	// remote.RemoteModel
	PullRequests     []*remote.PullRequest
	Issues           []*remote.Issue
	GithubCommitters []remote.Committer
}

// Create an enriched model by merging the local and GitHub model.
func NewEnrichedModel(local local.GitModel, github remote.RemoteModel) *EnrichedModel {
	return &EnrichedModel{
		// local.GitModel
		Commits:         local.Commits,
		Branches:        local.Branches,
		MainGraph:       local.MainGraph,
		BranchMatrix:    local.BranchMatrix,
		LocalCommitters: local.Committer,

		// remote.RemoteModel
		Name:             github.Name,
		URL:              github.URL,
		PullRequests:     github.PullRequests,
		Issues:           github.Issues,
		Owner:            github.Owner,
		GithubCommitters: github.Committers,
	}
}

func PopulateAuthors( //nolint: ireturn
	enriched *EnrichedModel,
	manualUsers ...struct{ email, login string },
) utils.Authors {
	authors := utils.NewAuthors()

	// Add manual committers.
	for _, m := range manualUsers {
		err := authors.Add(m.login, m.email)
		if err != nil {
			log.Fatalf("Error adding manual user: %v", err)
		}
	}

	// enriched model not available.
	if enriched == nil || enriched.GithubCommitters == nil {
		return authors
	}

	unavailableMap := make(map[string]struct{})    // map of unavailable committers.
	commitMap := make(map[string]remote.Committer) // map of commits to commitID.

	for _, committer := range enriched.GithubCommitters {
		commitMap[committer.CommitId] = committer
		if authors.Check(committer.Email) {
			continue
		}

		// Login is not always available.
		if committer.Login == "" {
			unavailableMap[committer.Email] = struct{}{}

			continue
		}

		err := authors.Add(committer.Login, committer.Email)
		if err != nil {
			log.Fatalf("Error adding committer: %v", err)
		}
	}

	for _, committer := range enriched.LocalCommitters {
		remoteCommitter, ok := commitMap[committer.CommitId]
		if !ok {
			continue
		}

		if committer.Email == remoteCommitter.Email {
			continue
		}

		if authors.Check(committer.Email) {
			continue
		}

		err := authors.Add(remoteCommitter.Login, committer.Email)
		if err != nil {
			log.Fatalf("Error adding committer: %v", err)
		}
	}

	unavailable := []string{}
	for u := range unavailableMap {
		unavailable = append(unavailable, u)
	}

	log.Println("Unavailable authors:", unavailable)

	return authors
}

// Find the current PR that the action is running on. Requires PR_NUMBER is set in env by action.
// See .github/workflows/git-gopher.yml for more info.
func (em *EnrichedModel) FindCurrentPR() (*remote.PullRequest, error) {
	// Pull request number from github action
	prNumberEnv := os.Getenv("PR_NUMBER")
	if prNumberEnv == "" {
		return nil, ErrPullRequestNumber
	}
	prNumber, err := strconv.Atoi(prNumberEnv)
	if err != nil {
		return nil, fmt.Errorf("could not atoi pr number: %w", err)
	}

	var targetPr *remote.PullRequest
	for _, pr := range em.PullRequests {
		if pr.Number == prNumber {
			targetPr = pr

			break
		}
	}

	if targetPr == nil {
		return nil, ErrFindPullRequest
	}

	return targetPr, nil
}
