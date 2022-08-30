package violation

import (
	"fmt"

	"github.com/Git-Gopher/go-gopher/markup"
)

func NewLinkedIssueViolation(
	pr markup.PR,
	current bool,
) *LinkedIssueViolation {
	violation := &LinkedIssueViolation{
		violation: violation{
			name:     "LinkedIssueViolation",
			severity: Suggestion,
			current:  current,
		},
		pr: pr,
	}
	violation.display = &display{violation}

	return violation
}

type LinkedIssueViolation struct {
	violation
	*display
	pr markup.PR
}

// Message implements Violation.
func (liv *LinkedIssueViolation) Message() string {
	return fmt.Sprintf("Pull request at %s does not contain a linked issue", liv.pr.Markdown())
}

// Suggestion implements Violation.
func (liv *LinkedIssueViolation) Suggestion() (string, error) {
	return "When appropriate (features, bug fixes), create and link an issue to the pull request. " +
		"This helps keep track of future developments and provides a development context for the pull request", nil
}
