package violation

import "fmt"

func NewBranchNameViolation(
	branchRef string,
	substring string,
) *BranchNameViolation {
	violation := &BranchNameViolation{
		display:   nil,
		branchRef: branchRef,
		substring: substring,
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
func (p *BranchNameViolation) Message() string {
	format := "Branch \"%s\" name is too inconsistent with other branch names"

	return fmt.Sprintf(format, p.branchRef)
}

// Name implements Violation.
func (*BranchNameViolation) Name() string {
	return "BranchNameViolation"
}

// Suggestion implements Violation.
func (p *BranchNameViolation) Suggestion() (string, error) {
	if p.substring == "" {
		return "", ErrViolationMethodNotExist
	}
	return fmt.Sprintf("All branch names should consistent with the substring \"%s\" ", p.substring), nil
}