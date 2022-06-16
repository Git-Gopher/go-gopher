package detector

import (
	"github.com/Git-Gopher/go-gopher/cache"
	"github.com/Git-Gopher/go-gopher/model/github"
	"github.com/Git-Gopher/go-gopher/violation"
)

type CommitCacheDetect func(current cache.Cache, cache []cache.Cache) (bool, violation.Violation, error)

type CommitCacheDetector struct {
	violated   int
	found      int
	total      int
	violations []violation.Violation

	detect CommitCacheDetect
}

// TODO: We should change this to the enriched model.
func (cd *CommitCacheDetector) Run(current cache.Cache, cache []cache.Cache) error {
	// Struct should be reset before each run, incase we are running it with a different model.
	cd.violated = 0
	cd.found = 0
	cd.total = 0
	cd.violations = make([]violation.Violation, 0)

	// Load cache internal to run?
	return ErrNotImplemented
}

func (cd *CommitCacheDetector) Run2(model *github.GithubModel) error {
	return nil
}

func (cd *CommitCacheDetector) Result() (int, int, int, []violation.Violation) {
	return cd.violated, cd.found, cd.total, cd.violations
}

func NewCommitCacheDetector(detect CommitCacheDetect) *CommitCacheDetector {
	return &CommitCacheDetector{
		violated:   0,
		found:      0,
		total:      0,
		violations: make([]violation.Violation, 0),
		detect:     detect,
	}
}
