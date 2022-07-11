package violation

import (
	"fmt"

	"github.com/Git-Gopher/go-gopher/model/github"
)

func NewBranchNameViolation(
	branchRef string,
	substring string,
	email string,
) *BranchNameViolation {
	violation := &BranchNameViolation{
		display:   nil,
		branchRef: branchRef,
		substring: substring,
		email:     email,
	}
	violation.display = &display{violation}

	return violation
}

// BranchNameViolation is violation when a branch name is inconsistent with others.
// from feature branches.
type BranchNameViolation struct {
	*display
	branchRef string
	substring string
	email     string
}

// FileLocation implements Violation.
func (*BranchNameViolation) FileLocation() (string, error) {
	return "", ErrViolationMethodNotExist
}

// LineLocation implements Violation.
func (*BranchNameViolation) LineLocation() (int, error) {
	return 0, ErrViolationMethodNotExist
}

// Message implements Violation.
func (bn *BranchNameViolation) Message() string {
	format := "Branch \"%s\" name is too inconsistent with other branch names"

	return fmt.Sprintf(format, bn.branchRef)
}

// Name implements Violation.
func (*BranchNameViolation) Name() string {
	return "BranchNameViolation"
}

// Suggestion implements Violation.
func (bn *BranchNameViolation) Suggestion() (string, error) {
	if bn.substring == "" {
		return "", ErrViolationMethodNotExist
	}

	return fmt.Sprintf("All branch names should consistent with the substring \"%s\" ", bn.substring), nil
}

// Author implements Violation.
func (bn *BranchNameViolation) Author() (*github.Author, error) {
	return nil, ErrViolationMethodNotExist
}

// Severity implements Violation.
func (bn *BranchNameViolation) Severity() Severity {
	return Suggestion
}

// Email implements Violation.
func (bn *BranchNameViolation) Email() string {
	return bn.email
}
