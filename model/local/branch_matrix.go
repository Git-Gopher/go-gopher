package local

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// Example of criss cross merge.
// This usually happen during hotfixes.
//
//		       3a4f5a6 -- 973b703 -- a34e5a1 (branch A)
//	        /        \ /
//
// 7c7bf85            X
//
//	        \        / \
//		        8f35f30 -- 3fd4180 -- 723181f (branch B)
type CrissCrossBranchInfo struct {
	Hash string
}

// BranchMatrix is an array of branch matrix.
// the matrix consists of all branches * all branches.
// e.g. A*B, A*C, B*C if branches A, B, C exist.
type BranchMatrix struct {
	A, B              *CrissCrossBranchInfo
	CrissCrossCommits []string
}

func CreateBranchMatrix(r *git.Repository, branchHeads []plumbing.Hash) ([]*BranchMatrix, error) {
	branchMatrix := []*BranchMatrix{}

	for _, a := range branchHeads {
		for _, b := range branchHeads {
			if a == b {
				continue
			}

			aCommits, err := r.CommitObject(a)
			if err != nil {
				return nil, fmt.Errorf("failed to get commit object a for %s: %w", a.String(), err)
			}

			bCommits, err := r.CommitObject(b)
			if err != nil {
				return nil, fmt.Errorf("failed to get commit object b for %s: %w", a.String(), err)
			}

			merge, err := aCommits.MergeBase(bCommits)
			if err != nil {
				return nil, fmt.Errorf("failed to get merge base for %s and %s: %w", a.String(), b.String(), err)
			}

			crissCrossCommits := []string{}
			for _, c := range merge {
				crissCrossCommits = append(crissCrossCommits, c.String())
			}

			branchMatrix = append(branchMatrix, &BranchMatrix{
				A: &CrissCrossBranchInfo{
					Hash: a.String(),
				},
				B: &CrissCrossBranchInfo{
					Hash: b.String(),
				},
				CrissCrossCommits: crissCrossCommits,
			})
		}
	}

	return branchMatrix, nil
}
