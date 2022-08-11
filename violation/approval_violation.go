package violation

import (
	"fmt"

	"github.com/Git-Gopher/go-gopher/markup"
)

func NewApprovalViolation(
	pr markup.PR,
) *ApprovalViolation {
	violation := &ApprovalViolation{
		violation: violation{
			name:     "PullRequestApprovalViolation",
			severity: Violated,
		},
		pr: pr,
	}
	violation.display = &display{violation}

	return violation
}

// from feature branches.
type ApprovalViolation struct {
	violation
	*display
	pr markup.PR
}

// Message implements Violation.
func (pra *ApprovalViolation) Message() string {
	return fmt.Sprintf("Pull request at %s was not reviewed before it was merged", pra.pr.Markdown())
}

// Suggestion implements Violation.
func (pra *ApprovalViolation) Suggestion() (string, error) {
	return "Try to approve pull requests before they are merged. " +
		"This gives your peers opportunity to look over your code and suggest improvements", nil
}
