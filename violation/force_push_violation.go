package violation

import (
	"fmt"
	"strings"
	"time"

	"github.com/Git-Gopher/go-gopher/markup"
)

func NewForcePushViolation(
	lostCommits []markup.Commit,
	email string,
	time time.Time,
) *ForcePushViolation {
	violation := &ForcePushViolation{
		violation: violation{
			name:     "ForcePushViolation",
			email:    email,
			time:     time,
			severity: Violated,
		},
		lostCommits: lostCommits,
	}
	violation.display = &display{violation}

	return violation
}

// Force push violation occurs whenever a branch is force pushed to, losing a series of commits
// from feature branches.
type ForcePushViolation struct {
	violation
	*display
	lostCommits []markup.Commit
}

// Message implements Violation.
func (fpv *ForcePushViolation) Message() string {
	format := "The following commits have been lost as result of a force push: %s"
	commits := make([]string, len(fpv.lostCommits))
	for i, commit := range fpv.lostCommits {
		commits[i] = commit.Markdown()
	}

	return fmt.Sprintf(format, strings.Join(commits, ",\n"))
}

// Suggestion implements Violation.
func (fpv *ForcePushViolation) Suggestion() (string, error) {
	format := "Restore the following commits to restore the work lost on the branch:\n%s"
	commits := make([]string, len(fpv.lostCommits))
	for i, commit := range fpv.lostCommits {
		commits[i] = commit.Markdown()
	}

	return fmt.Sprintf(format, strings.Join(commits, ",")), nil
}

// Current implements Violation.
func (fpv *ForcePushViolation) Current() bool {
	return true
}
