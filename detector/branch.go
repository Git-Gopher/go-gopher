package detector

import (
	"fmt"
	"time"

	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/violation"
)

type BranchDetect func(branch *local.Branch) (bool, violation.Violation, error)

// BranchDetector is used to run a detector on each branch metadata.
type BranchDetector struct {
	violated   int
	found      int
	total      int
	violations []violation.Violation

	detect BranchDetect
}

func (bd *BranchDetector) Run(model *enriched.EnrichedModel) error {
	bd.violated = 0
	bd.found = 0
	bd.total = 0
	bd.violations = make([]violation.Violation, 0)

	for _, b := range model.Branches {
		b := b
		detected, violation, err := bd.detect(&b)
		if err != nil {
			return fmt.Errorf("Error detecting stale branch: %w", err)
		}
		if err != nil {
			return err
		}
		if detected {
			bd.found++
		}
		if violation != nil {
			bd.violations = append(bd.violations, violation)
		}
		bd.total++
	}

	return nil
}

func (b *BranchDetector) Result() (int, int, int, []violation.Violation) {
	return b.violated, b.found, b.total, b.violations
}

func NewBranchDetector(detect BranchDetect) *BranchDetector {
	return &BranchDetector{
		violated:   0,
		found:      0,
		total:      0,
		violations: make([]violation.Violation, 0),
		detect:     detect,
	}
}

// GithubWorklow: Branches are considered stale after three months.
func StaleBranchDetect() BranchDetect {
	StaleBranchTime := time.Hour * 24 * 30

	return func(branch *local.Branch) (bool, violation.Violation, error) {
		if time.Since(branch.Head.Committer.When) > StaleBranchTime {
			return true, violation.NewStaleBranchViolation(branch.Name, StaleBranchTime), nil
		}

		return false, nil, nil
	}
}
