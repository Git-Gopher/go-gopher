package violation

import (
	"fmt"
	"time"

	"github.com/Git-Gopher/go-gopher/markup"
)

func NewBranchNameViolation(
	branchRef markup.Branch,
	substring string,
	email string,
	time time.Time,
) *BranchNameViolation {
	violation := &BranchNameViolation{
		violation: violation{
			name:     "BranchNameViolation",
			email:    email,
			time:     time,
			severity: Suggestion,
		},
		branchRef: branchRef,
		substring: substring,
	}
	violation.display = &display{violation}

	return violation
}

// BranchNameViolation is violation when a branch name is inconsistent with others.
// from feature branches.
type BranchNameViolation struct {
	violation
	*display
	branchRef markup.Branch
	substring string
}

// Message implements Violation.
func (bn *BranchNameViolation) Message() string {
	format := "Branch %s name is too inconsistent with other branch names"

	return fmt.Sprintf(format, bn.branchRef.Markdown())
}

// Suggestion implements Violation.
func (bn *BranchNameViolation) Suggestion() (string, error) {
	if bn.substring == "" {
		return "", ErrViolationMethodNotExist
	}

	return fmt.Sprintf("All branch names should consistent with the substring \"%s\" ", bn.substring), nil
}
