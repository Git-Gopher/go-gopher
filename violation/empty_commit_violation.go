package violation

import (
	"fmt"
	"time"

	"github.com/Git-Gopher/go-gopher/markup"
)

func NewEmptyCommitViolation(
	commit markup.Commit,
	email string,
	time time.Time,
) *EmptyCommitViolation {
	violation := &EmptyCommitViolation{
		violation: violation{
			name:     "EmptyCommitViolation",
			email:    email,
			time:     time,
			severity: Violated,
		},
		commit: commit,
	}
	violation.display = &display{violation}

	return violation
}

// Empty commit violation is a violation when a commit is made with no content.
// Students may aim to inflate their commit count artificially by making empty commits.
type EmptyCommitViolation struct {
	violation
	*display
	commit markup.Commit
}

// Message implements Violation.
func (ecv *EmptyCommitViolation) Message() string {
	format := "Commit %s is empty but has been committed to the project"

	return fmt.Sprintf(format, ecv.commit.Markdown())
}

// Suggestion implements Violation.
func (ecv *EmptyCommitViolation) Suggestion() (string, error) {
	return "Try not make empty commits to the git history " +
		"as it makes it seem like you are forging or padding your version history", nil
}
