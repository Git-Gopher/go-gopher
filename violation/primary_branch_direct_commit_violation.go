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
) *PrimaryBranchDirectCommitViolation {
	violation := &PrimaryBranchDirectCommitViolation{
		violation: violation{
			name:     "PrimaryBranchDirectCommitViolation",
			email:    email,
			time:     time,
			severity: Violated,
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
func (p *PrimaryBranchDirectCommitViolation) Message() string {
	format := "Commit \"%s\" has been directly committed to the primary branch \"%s\""

	return fmt.Sprintf(format, p.commitHash.Link(), p.primaryBranch.Link())
}

// Suggestion implements Violation.
func (p *PrimaryBranchDirectCommitViolation) Suggestion() (string, error) {
	return fmt.Sprintf("All commits should be merged in to the branch \"%s\" ", p.primaryBranch.Link()), nil
}
