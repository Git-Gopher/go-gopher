package violation

import (
	"fmt"
	"time"

	"github.com/Git-Gopher/go-gopher/markup"
)

func NewApprovalViolation(
	pr markup.PR,
	current bool,
	time time.Time,
) *ApprovalViolation {
	violation := &ApprovalViolation{
		violation: violation{
			name:     "PullRequestApprovalViolation",
			severity: Violated,
			time:     time,
			current:  current,
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
func (av *ApprovalViolation) Message() string {
	return fmt.Sprintf("Pull request at %s did not receive review approval before it was merged", av.pr.Markdown())
}

// Suggestion implements Violation.
func (av *ApprovalViolation) Suggestion() (string, error) {
	return "Ensure that you are reviewing pull requests before they get merged. Reviews can be added to a pull request " +
			"by checking (Files changed > Review changes) on the pull request. Pull requests should receive at least one " +
			"approval before they are merged, this gives your peers opportunity to look over your code and suggest improvements",
		nil
}
