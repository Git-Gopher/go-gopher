package violation

import (
	"time"

	"github.com/Git-Gopher/go-gopher/markup"
)

func NewStaleCommitViolation(commit markup.Commit, message string, email string, time time.Time) *StaleCommitViolation {
	common := &StaleCommitViolation{
		violation: violation{
			name:     "StaleCommitViolation",
			email:    email,
			time:     time,
			severity: Violated,
		},
		commit:  commit,
		message: message,
	}
	common.display = &display{common}

	return common
}

// Example violation.
type StaleCommitViolation struct {
	violation
	*display
	commit  markup.Commit
	message string
}

// Message implements Violation.
func (sc *StaleCommitViolation) Message() string {
	return sc.message
}
