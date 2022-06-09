package enriched

import (
	"github.com/Git-Gopher/go-gopher/model/github"
	"github.com/Git-Gopher/go-gopher/model/local"
)

type EnrichedModel struct {
}

func NewEnriched(local *local.GitModel, github *github.GithubModel) (*EnrichedModel, error) {
	return nil, nil
}
