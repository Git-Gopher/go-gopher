package violation

import (
	"fmt"

	"github.com/Git-Gopher/go-gopher/model/remote"
	"github.com/Git-Gopher/go-gopher/utils"
)

func NewPrimaryBranchDirectCommitViolation(
	primaryBranch utils.Branch,
	commitHash utils.Commit,
	parentHashes []utils.Commit,
	email string,
) *PrimaryBranchDirectCommitViolation {
	violation := &PrimaryBranchDirectCommitViolation{
		display:       nil,
		parentHashes:  parentHashes,
		primaryBranch: primaryBranch,
		commitHash:    commitHash,
		email:         email,
	}
	violation.display = &display{violation}

	return violation
}

// PrimaryBranchDirectCommitViolation is violation when a commit is done directly to primary branch without merging
// from feature branches.
type PrimaryBranchDirectCommitViolation struct {
	*display
	primaryBranch utils.Branch
	parentHashes  []utils.Commit
	commitHash    utils.Commit
	email         string
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

	return fmt.Sprintf(format, p.commitHash.Link(), p.primaryBranch.Link())
}

// Name implements Violation.
func (*PrimaryBranchDirectCommitViolation) Name() string {
	return "PrimaryBranchDirectCommitViolation"
}

// Suggestion implements Violation.
func (p *PrimaryBranchDirectCommitViolation) Suggestion() (string, error) {
	return fmt.Sprintf("All commits should be merged in to the branch \"%s\" ", p.primaryBranch.Link()), nil
}

// Author implements Violation.
func (p *PrimaryBranchDirectCommitViolation) Author() (*remote.Author, error) {
	return nil, ErrViolationMethodNotExist
}

// Severity implements Violation.
func (p *PrimaryBranchDirectCommitViolation) Severity() Severity {
	return Violated
}

// Severity implements Violation.
func (p *PrimaryBranchDirectCommitViolation) Email() string {
	return p.email
}
