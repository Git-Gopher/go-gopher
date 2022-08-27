package violation

import (
	"fmt"
	"time"

	"github.com/Git-Gopher/go-gopher/markup"
)

func NewUnresolvedMergeViolation(
	line markup.Line,
	email string,
	time time.Time,
) *UnresolvedMergeViolation {
	violation := &UnresolvedMergeViolation{
		violation: violation{
			name:     "UnresolvedMergeViolation",
			email:    email,
			time:     time,
			severity: Violated,
		},
		line: line,
	}
	violation.display = &display{violation}

	return violation
}

// Violation when unresolved merge conflicts are merged violation.
type UnresolvedMergeViolation struct {
	violation
	*display
	line markup.Line
}

// Message implements Violation.
func (umv *UnresolvedMergeViolation) Message() string {
	return fmt.Sprintf("Unresolved merge conflicts at line %s", umv.line.Markdown())
}

// Suggestion implements Violation.
func (umv *UnresolvedMergeViolation) Suggestion() (string, error) {
	return "Resolve merge conflicts before committing to a branch " +
		"as unresolved conflicts cause the project to enter a broken state", nil
}

// Current implements Violation.
func (umv *UnresolvedMergeViolation) Current() bool {
	return true
}
