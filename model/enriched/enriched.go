package enriched

import (
	"github.com/Git-Gopher/go-gopher/model/github"
	"github.com/Git-Gopher/go-gopher/model/local"
)

type EnrichedModel struct {
	Commits      []local.Commit
	Branches     []local.Branch
	PullRequests []*github.PullRequest
	Issues       []*github.Issue
}

// Create an enriched model by merging the local and GitHub model.
func NewEnrichedModel(local local.GitModel, github github.GithubModel) *EnrichedModel {
	return &EnrichedModel{
		Commits:      local.Commits,
		Branches:     local.Branches,
		PullRequests: github.PullRequests,
		Issues:       github.Issues,
	}
}
