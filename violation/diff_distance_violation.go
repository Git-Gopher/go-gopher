package violation

import (
	"fmt"
	"time"

	"github.com/Git-Gopher/go-gopher/markup"
)

func NewExtremeDiffDistanceViolation(
	commit markup.Commit,
	email string,
	time time.Time,
	current bool,
) *DiffDistanceViolation {
	common := &DiffDistanceViolation{
		violation: violation{
			name:     "ExtremeDiffDistanceViolation",
			email:    email,
			time:     time,
			severity: Suggestion,
			current:  current,
		},
		commit: commit,
	}
	common.display = &display{common}

	return common
}

func NewMildDiffDistanceViolation(
	commit markup.Commit,
	email string,
	time time.Time,
	current bool,
) *DiffDistanceViolation {
	common := &DiffDistanceViolation{
		violation: violation{
			name:     "MildDiffDistanceViolation",
			email:    email,
			time:     time,
			severity: Suggestion,
			current:  current,
		},
		commit: commit,
	}
	common.display = &display{common}

	return common
}

// Example violation.
type DiffDistanceViolation struct {
	violation
	*display
	commit markup.Commit
}

// Message implements Violation.
func (ddv *DiffDistanceViolation) Message() string {
	return fmt.Sprintf("Fragmented commit found at %s. This commit might contain multiple changes", ddv.commit.Markdown())
}

func (ddv *DiffDistanceViolation) Suggestion() (string, error) {
	return "Commits should aim to tackle one problem at a time. Try to break down your tasks into achievable " +
		"subtasks and create a commit when you feel that the subtask has been completed. This way you may revert " +
		"to any stage of your task at any point without losing the progress of multiple changes", nil
}
