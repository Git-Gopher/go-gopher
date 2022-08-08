package violation

import (
	"fmt"
)

func NewLinkedIssueViolation(
	url string,
) *LinkedIssueViolation {
	violation := &LinkedIssueViolation{
		violation: violation{
			name:     "LinkedIssueViolation",
			severity: Violated,
		},
		url: url,
	}
	violation.display = &display{violation}

	return violation
}

type LinkedIssueViolation struct {
	violation
	*display
	url string
}

// Message implements Violation.
func (liv *LinkedIssueViolation) Message() string {
	return fmt.Sprintf("Pull request at %s does not contain a linked issue", liv.url)
}

// Suggestion implements Violation.
func (liv *LinkedIssueViolation) Suggestion() (string, error) {
	return `When appropriate (features, bug fixes), create and link an issue to the pull request. 
	This helps keep track of future developments and provides a development context for the pull request`, nil
}
