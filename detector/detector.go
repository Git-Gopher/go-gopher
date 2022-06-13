package detector

import (
	"github.com/Git-Gopher/go-gopher/model/github"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/violation"
)

type Detector interface {
	Result() (violated, count, total int, violations []violation.Violation)
	// TODO: We should change this to the enriched model
	Run(model *local.GitModel) error
	Run2(model *github.GithubModel) error
}
