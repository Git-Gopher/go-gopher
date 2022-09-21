package detector

import (
	"errors"

	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/violation"
)

var ErrDevelopBranchModelNil = errors.New("develop branch model is nil")

// DevelopBranch is a detector that checks if release and develop branch exist.
type DevelopBranch struct {
	name       string
	violated   int // non feature branches aka develop/release etc. (does not account default branch)
	found      int // total feature branches
	total      int // total branches
	violations []violation.Violation
}

// NewDevelopBranch creates a new develop branch detector.
func NewDevelopBranch(name string) *CherryPickDetector {
	return &CherryPickDetector{
		name:       name,
		violated:   0,
		found:      0,
		total:      0,
		violations: make([]violation.Violation, 0),
	}
}

func (db *DevelopBranch) Run(em *enriched.EnrichedModel) error {
	if em == nil {
		return ErrDevelopBranchModelNil
	}

	db.violated = 0
	db.found = 0
	db.total = 0
	db.violations = make([]violation.Violation, 0)

	if em.ReleaseGraph != nil && em.MainGraph != nil {
		db.found = 1
	}

	db.total = 1

	return nil
}

func (db *DevelopBranch) Result() (int, int, int, []violation.Violation) {
	return db.violated, db.found, db.total, db.violations
}

func (db *DevelopBranch) Name() string {
	return db.name
}
