package detector

import (
	"fmt"

	"github.com/Git-Gopher/go-gopher/cache"
	"github.com/Git-Gopher/go-gopher/model/github"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/violation"
)

var ErrNotImplemented = fmt.Errorf("Not implemented")

type Detector interface {
	// TODO: We should change this to the enriched model
	Run(model *local.GitModel) error
	Run2(model *github.GithubModel) error
	Result() (violated, count, total int, violations []violation.Violation)
}

type CacheDetector interface {
	Run(current *cache.Cache, cache []*cache.Cache) error
	Result() (violated, count, total int, violations []violation.Violation)
}
