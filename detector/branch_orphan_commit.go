package detector

import (
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/violation"
)

// TODO: Github Workflow: If a branch remains unmerged then it is considered to be stale.
// Rules:
// - If the head commit of the branch is more than X days/weeks old then branch is old (unused)
// - If branch has fallen behind primary branch by a certain amount then it has been failed to be maintained.
// Use CreateIndependentCommits.
func NewOrphanCommitDetector() BranchDetect {
	return func(branch *local.Branch) (bool, violation.Violation, error) {
		// Run the detector on each branch and look for orphan commits

		return false, nil, nil
	}
}
