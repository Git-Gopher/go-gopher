package detector

import (
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/violation"
)

type BranchMatrixDetect func(branchMatrix *local.BranchMatrix) (bool, violation.Violation, error)

// FeatureBranchDetector is a detector that detects multiple remote branches (not deleted).
// And check if the branch is a feature branch or a main/develop branch.
type BranchMatrixDetector struct {
	name       string
	violated   int // no violations
	found      int // total branches with cross-merge
	total      int // total branches * total branches
	violations []violation.Violation

	detect BranchMatrixDetect
}

func NewBranchMatrixDetector(name string, detect BranchMatrixDetect) *BranchMatrixDetector {
	return &BranchMatrixDetector{
		name:       name,
		total:      0,
		found:      0,
		detect:     detect,
		violations: make([]violation.Violation, 0),
	}
}

func (cc *BranchMatrixDetector) Run(model *enriched.EnrichedModel) error {
	for _, b := range model.BranchMatrix {
		b := b
		detected, violation, err := cc.detect(b)
		cc.total++
		if err != nil {
			return err
		}
		if violation != nil {
			cc.violations = append(cc.violations, violation)
		}
		if detected {
			cc.found++
		}
	}

	return nil
}

func (cc *BranchMatrixDetector) Result() (violated int, count int, total int, violations []violation.Violation) {
	return cc.violated, cc.found, cc.total, cc.violations
}

func (cc *BranchMatrixDetector) Name() string {
	return cc.name
}

// CrissCrossMergeDetect to find criss cross merges
// Example of criss cross merge.
// This usually happen during hotfixes.
//
//	         3a4f5a6 -- 973b703 -- a34e5a1 (branch A)
//	       /        \ /
//
//	7c7bf85          X
//
//	       \        / \
//	         8f35f30 -- 3fd4180 -- 723181f (branch B)
func CrissCrossMergeDetect() (string, BranchMatrixDetect) {
	return "CrissCrossMergeDetect", func(branchMatrix *local.BranchMatrix) (bool, violation.Violation, error) {
		if len(branchMatrix.CrissCrossCommits) >= 2 {
			return true, nil, nil
		}

		return false, nil, nil
	}
}
