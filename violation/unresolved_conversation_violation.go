package violation

import (
	"fmt"
	"time"

	"github.com/Git-Gopher/go-gopher/markup"
)

func NewUnresolvedConversationViolation(
	pr markup.PR,
	current bool,
	time time.Time,
) *UnresolvedConversationViolation {
	violation := &UnresolvedConversationViolation{
		violation: violation{
			name:     "UnresolvedConversationViolation",
			severity: Violated,
			time:     time,
			current:  current,
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
	return fmt.Sprintf("Merged pull request at %s contains unresolved conversation threads", ucv.pr.Markdown())
}

// Suggestion implements Violation.
func (ucv *UnresolvedConversationViolation) Suggestion() (string, error) {
	return "Resolve conversations before pull requests are merged. " +
		"This indicates that all discussions for the pull request have been resolved and " +
		"puts the entire team is on the same page with the progress of your project " +
		"without having to double check with your peers", nil
}
