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
	format := "Branch %s name might inconsistent with other branch names"

	return fmt.Sprintf(format, bn.branchRef.Markdown())
}

// Suggestion implements Violation.
func (bn *BranchNameViolation) Suggestion() (string, error) {
	if bn.substring == "" {
		return "", ErrViolationMethodNotExist
	}

	return fmt.Sprintf("Try to group together branch names by using a group token prefix that indicates "+
		"the type of change that the branch contains. For example \"fix/\" or \"feature/\" are "+
		"good group tokens for branches containing fixes and features respectively."+
		"The current longest prefix your branches have is \"%s\"", bn.substring), nil
}
