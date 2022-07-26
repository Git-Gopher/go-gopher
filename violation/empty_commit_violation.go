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
			severity: Suggestion,
		},
		commit: commit,
	}
	violation.display = &display{violation}

	return violation
}

// BinaryViolation is violation when a branch name is inconsistent with others.
// from feature branches.
type EmptyCommitViolation struct {
	violation
	*display
	commit markup.Commit
}

// Message implements Violation.
func (ecv *EmptyCommitViolation) Message() string {
	format := "Empty commit %s has been committed to the project"

	return fmt.Sprintf(format, ecv.commit.Markdown())
}

// Suggestion implements Violation.
func (ecv *EmptyCommitViolation) Suggestion() (string, error) {
	return `Try not make empty commits to the project 
		as it makes it seem like you are forging or padding your version history`, nil
}