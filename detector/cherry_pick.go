package detector

import (
	"errors"

	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/violation"
)

var (
	ErrCherryPickModelNil        = errors.New("cherry pick model is nil")
	ErrCherryPickNoReleaseBranch = errors.New("cherry pick detector: no release branch")
	ErrCherryPickNoDevBranch     = errors.New("cherry pick detector: no dev branch")
)

// CherryPickDetector is a detector that counts the number of cherry picked commits.
// found / total = number of cherry picked commits found with compariable patch id.
// only unconflicted cherry pick can be picked up.
type CherryPickDetector struct {
	name       string
	violated   int // non feature branches aka develop/release etc. (does not account default branch)
	found      int // total feature branches
	total      int // total branches
	violations []violation.Violation
}

// NewCherryPickDetector creates a new cherry pick detector.
func NewCherryPickDetector(name string) *CherryPickDetector {
	return &CherryPickDetector{
		name:       name,
		violated:   0,
		found:      0,
		total:      0,
		violations: make([]violation.Violation, 0),
	}
}

func (cp *CherryPickDetector) Run(em *enriched.EnrichedModel) error {
	if em == nil {
		return ErrCherryPickModelNil
	}

	patchMap := make(map[string]string)

	for _, commit := range em.Commits {
		if commit.PatchID == nil {
			continue
		}

		cp.total++

		if hash, ok := patchMap[*commit.PatchID]; !ok {
			patchMap[*commit.PatchID] = commit.Hash.HexString()
		} else if hash != commit.Hash.HexString() {
			cp.found++
		}
	}

	return nil
}

func (cp *CherryPickDetector) Result() (int, int, int, []violation.Violation) {
	return cp.violated, cp.found, cp.total, cp.violations
}

func (cp *CherryPickDetector) Name() string {
	return cp.name
}
