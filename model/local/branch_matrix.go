package local

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
