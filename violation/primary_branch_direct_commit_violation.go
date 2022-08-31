package violation

import (
	"fmt"
	"time"

	"github.com/Git-Gopher/go-gopher/markup"
)

func NewPrimaryBranchDirectCommitViolation(
	primaryBranch markup.Branch,
	commitHash markup.Commit,
	parentHashes []markup.Commit,
	email string,
	time time.Time,
	current bool,
) *PrimaryBranchDirectCommitViolation {
	violation := &PrimaryBranchDirectCommitViolation{
		violation: violation{
			name:     "PrimaryBranchDirectCommitViolation",
			email:    email,
			time:     time,
			severity: Violated,
			current:  current,
		},
		parentHashes:  parentHashes,
		primaryBranch: primaryBranch,
		commitHash:    commitHash,
	}
	violation.display = &display{violation}

	return violation
}

// PrimaryBranchDirectCommitViolation is violation when a commit is done directly to primary branch without merging
// from feature branches.
type PrimaryBranchDirectCommitViolation struct {
	violation
	*display
	primaryBranch markup.Branch
	parentHashes  []markup.Commit
	commitHash    markup.Commit
}

// Message implements Violation.
func (pbdvc *PrimaryBranchDirectCommitViolation) Message() string {
	format := "Commit %s has been directly committed to the primary branch %s"

	return fmt.Sprintf(format, pbdvc.commitHash.Markdown(), pbdvc.primaryBranch.Markdown())
}

// Suggestion implements Violation.
func (pbdcv *PrimaryBranchDirectCommitViolation) Suggestion() (string, error) {
	return fmt.Sprintf("All commits should be merged in to the branch %s via a pull request, "+
		"instead of directly committing to the primary branch. This method helps keep track of your development history "+
		"and is a fundamental technique of using Github Flow",
		pbdcv.primaryBranch.Markdown()), nil
}
