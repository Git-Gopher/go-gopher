package enriched

import (
	"log"

	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/model/remote"
	"github.com/Git-Gopher/go-gopher/utils"
)

type EnrichedModel struct {
	// local.GitModel
	Owner        string                `json:"owner"`
	Name         string                `json:"name"`
	URL          string                `json:"url"`
	Commits      []local.Commit        `json:"commits"`
	Branches     []local.Branch        `json:"branches"`
	MainGraph    *local.BranchGraph    `json:"mainGraph"`    // Graph representation of commits in the main branch
	BranchMatrix []*local.BranchMatrix `json:"branchMatrix"` // Matrix representation by comparing branches

	// remote.RemoteModel
	PullRequests []*remote.PullRequest `json:"pullRequests"`
	Issues       []*remote.Issue       `json:"issues"`
	Committers   []remote.Committer    `json:"committers"`
}

// Create an enriched model by merging the local and GitHub model.
func NewEnrichedModel(local local.GitModel, github remote.RemoteModel) *EnrichedModel {
	return &EnrichedModel{
		// local.GitModel
		Commits:      local.Commits,
		Branches:     local.Branches,
		MainGraph:    local.MainGraph,
		BranchMatrix: local.BranchMatrix,

		// remote.RemoteModel
		Name:         github.Name,
		URL:          github.URL,
		PullRequests: github.PullRequests,
		Issues:       github.Issues,
		Owner:        github.Owner,
		Committers:   github.Committers,
	}
}

// nolint:ireturn
func PopulateAuthors(enriched *EnrichedModel, manualUsers ...struct{ email, login string }) utils.Authors {
	if enriched == nil || enriched.Committers == nil {
		return utils.NewAuthors()
	}

	authors := utils.NewAuthors()

	for _, m := range manualUsers {
		err := authors.Add(m.login, m.email)
		if err != nil {
			log.Fatalf("Error adding manual user: %v", err)
		}
	}

	unavailableMap := make(map[string]struct{})

	for _, committer := range enriched.Committers {
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

	unavailable := []string{}
	for u := range unavailableMap {
		unavailable = append(unavailable, u)
	}

	log.Println("Unavailable authors:", unavailable)

	return authors
}
