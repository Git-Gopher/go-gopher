package detector

import "log"

// Example of criss cross merge.
// This usually happen during hotfixes.
//
//          3a4f5a6 -- 973b703 -- a34e5a1 (branch A)
//        /        \ /
// 7c7bf85          X
//        \        / \
//          8f35f30 -- 3fd4180 -- 723181f (branch B)

type CrissCrossBranchInfo struct {
	Hash string
}

// BranchMatrixModel is an array of branch matrix.
// the matrix consists of all branches * all branches.
// e.g. A*B, A*C, B*C if branches A, B, C exist.
type BranchMatrixModel struct {
	A, B              *CrissCrossBranchInfo
	CrissCrossCommits []string
}

type Violation struct {
	Info string
}

type BranchMatrixDetect func(branchMatrix *BranchMatrixModel) (bool, *Violation, error)

// FeatureBranchDetector is a detector that detects multiple remote branches (not deleted).
// And check if the branch is a feature branch or a main/develop branch.
type BranchMatrixDetector struct {
	violated int // no violations
	found    int // total branches with cross-merge
	total    int // total branches * total branches

	detect BranchMatrixDetect
}

// TODO: Move into GitModel

func (cc *BranchMatrixDetector) Run(branchMatrix []BranchMatrixModel) error {
	for _, b := range branchMatrix {
		b := b
		detected, violation, err := cc.detect(&b)
		cc.total++
		if err != nil {
			return err
		}
		if violation != nil {
			// TODO: implement violation log handler
			log.Println(violation.Info)
		}
		if detected {
			cc.found++
		}
	}

	return nil
}

func (cc *BranchMatrixDetector) Result() (violated, count, total int) {
	return cc.violated, cc.found, cc.total
}

func NewBranchMatrixDetector(detect BranchMatrixDetect) *BranchMatrixDetector {
	return &BranchMatrixDetector{
		total:  0,
		found:  0,
		detect: detect,
	}
}

func NewCrissCrossMergeDetect() BranchMatrixDetect {
	return func(branchMatrix *BranchMatrixModel) (bool, *Violation, error) {
		if len(branchMatrix.CrissCrossCommits) >= 2 {
			return true, nil, nil
		}

		return false, nil, nil
	}
}
