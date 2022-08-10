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
func (um *UnresolvedMergeViolation) Message() string {
	return fmt.Sprintf("Unresolved merge conflicts on %s", um.line.Markdown())
}

// Suggestion implements Violation.
func (um *UnresolvedMergeViolation) Suggestion() (string, error) {
	return `Resolve conflicts before committing. 
		This ensures that the project is not left in a broken state for others`, nil
}
