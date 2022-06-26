package violation

import (
	"fmt"

	"github.com/Git-Gopher/go-gopher/model/github"
)

func NewPrimaryBranchDirectCommitViolation(
	primaryBranch string,
	commitHash string,
	parentHashes []string,
) *PrimaryBranchDirectCommitViolation {
	violation := &PrimaryBranchDirectCommitViolation{
		display:       nil,
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
	*display
	parentHashes  []string
	primaryBranch string
	commitHash    string
}

// FileLocation implements Violation.
func (*PrimaryBranchDirectCommitViolation) FileLocation() (string, error) {
	return "", ErrViolationMethodNotExist
}

// LineLocation implements Violation.
func (*PrimaryBranchDirectCommitViolation) LineLocation() (int, error) {
	return 0, ErrViolationMethodNotExist
}

// Message implements Violation.
func (p *PrimaryBranchDirectCommitViolation) Message() string {
	format := "Commit \"%s\" has been directly committed to the primary branch \"%s\""

	return fmt.Sprintf(format, p.commitHash, p.primaryBranch)
}

// Name implements Violation.
func (*PrimaryBranchDirectCommitViolation) Name() string {
	return "PrimaryBranchDirectCommitViolation"
}

// Suggestion implements Violation.
func (p *PrimaryBranchDirectCommitViolation) Suggestion() (string, error) {
	return fmt.Sprintf("All commits should be merged in to the branch \"%s\" ", p.primaryBranch), nil
}

// Author implements Violation.
func (p *PrimaryBranchDirectCommitViolation) Author() (*github.Author, error) {
	return nil, ErrViolationMethodNotExist
}

// Severity implements Violation.
func (p *PrimaryBranchDirectCommitViolation) Severity() Severity {
	return Violated
}
