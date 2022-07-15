package detector

import (
	"fmt"

	"github.com/Git-Gopher/go-gopher/cache"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/violation"
)

var ErrNotImplemented = fmt.Errorf("Not implemented")

// common - common variables that are shared with all detectors
type common struct {
	owner string
	repo  string
}

type Detector interface {
	Run(model *enriched.EnrichedModel) error
	Result() (violated, count, total int, violations []violation.Violation)
}

type CacheDetector interface {
	Run(owner string, repo string, email string, current *cache.Cache, cache []*cache.Cache) error
	Result() (violated, count, total int, violations []violation.Violation)
}
