package enriched

import (
	"github.com/Git-Gopher/go-gopher/model/github"
	"github.com/Git-Gopher/go-gopher/model/local"
)

type EnrichedModel struct {
	// local.GitModel
	Commits      []local.Commit
	Branches     []local.Branch
	MainGraph    *local.BranchGraph    // Graph representation of commits in the main branch
	BranchMatrix []*local.BranchMatrix // Matrix representation by comparing branches

	// github.GithubModel
	PullRequests []*github.PullRequest
	Issues       []*github.Issue
}

// Create an enriched model by merging the local and GitHub model.
func NewEnrichedModel(local local.GitModel, github github.GithubModel) *EnrichedModel {
	return &EnrichedModel{
		// local.GitModel
		Commits:      local.Commits,
		Branches:     local.Branches,
		MainGraph:    local.MainGraph,
		BranchMatrix: local.BranchMatrix,

		// github.GithubModel
		PullRequests: github.PullRequests,
		Issues:       github.Issues,
	}
}
