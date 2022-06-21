package detector

import (
	"fmt"
	"time"

	"github.com/Git-Gopher/go-gopher/model/github"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/violation"
)

var StaleBranchTime = time.Hour * 24 * 30

type BranchDetect func(branch *local.Branch) (bool, violation.Violation, error)

// BranchDetector is used to run a detector on each branch metadata.
type BranchDetector struct {
	violated   int
	found      int
	total      int
	violations []violation.Violation

	detect BranchDetect
}

// TODO: We should change this to the enriched model.
func (bd *BranchDetector) Run(model *local.GitModel) error {
	bd.violated = 0
	bd.found = 0
	bd.total = 0
	bd.violations = make([]violation.Violation, 0)

	for _, b := range model.Branches {
		detected, violation, err := bd.detect(&b)
		if err != nil {
			return fmt.Errorf("Error detecting stale branch: %v", err)
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

func (db *BranchDetector) Run2(model *github.GithubModel) error {
	return ErrNotImplemented
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
	return func(branch *local.Branch) (bool, violation.Violation, error) {
		if time.Since(branch.Head.Committer.When) > StaleBranchTime {
			return true, violation.NewStaleBranchViolation(branch.Name, StaleBranchTime), nil
		}

		return false, nil, nil
	}
}
