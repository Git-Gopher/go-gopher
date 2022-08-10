package violation

import (
	"fmt"

	"github.com/Git-Gopher/go-gopher/markup"
)

func NewUnresolvedConversationViolation(
	pr markup.PR,
) *UnresolvedConversationViolation {
	violation := &UnresolvedConversationViolation{
		violation: violation{
			name:     "UnresolvedConversationViolation",
			severity: Violated,
		},
		pr: pr,
	}
	violation.display = &display{violation}

	return violation
}

type UnresolvedConversationViolation struct {
	violation
	*display
	pr markup.PR
}

// Message implements Violation.
func (ucv *UnresolvedConversationViolation) Message() string {
	return fmt.Sprintf("Pull request at %s contains unresolved conversation threads", ucv.pr.Markdown())
}

// Suggestion implements Violation.
func (ucv *UnresolvedConversationViolation) Suggestion() (string, error) {
	return `Try to resolve conversations before they are merged. 
	This ensures that everyone on your team is on the same page with the progress of a particular pull request
	without having to check your chat history`, nil
}
