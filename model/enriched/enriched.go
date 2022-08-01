package enriched

import (
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/model/remote"
	"github.com/Git-Gopher/go-gopher/utils"
	log "github.com/sirupsen/logrus"
)

type EnrichedModel struct {
	// local.GitModel
	Owner           string
	Name            string
	URL             string
	Commits         []local.Commit
	Branches        []local.Branch
	MainGraph       *local.BranchGraph    // Graph representation of commits in the main branch
	BranchMatrix    []*local.BranchMatrix // Matrix representation by comparing branches
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

// nolint:ireturn
func PopulateAuthors(enriched *EnrichedModel, manualUsers ...struct{ email, login string }) utils.Authors {
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
