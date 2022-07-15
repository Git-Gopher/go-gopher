package violation

import (
	"github.com/Git-Gopher/go-gopher/model/remote"
	"github.com/Git-Gopher/go-gopher/utils"
)

func NewStaleCommitViolation(commit utils.Commit, message string) *StaleCommitViolation {
	common := &StaleCommitViolation{
		display: nil,
		commit:  commit,
		message: message,
	}
	common.display = &display{common}

	return common
}

// Example violation.
type StaleCommitViolation struct {
	*display
	commit  utils.Commit
	message string
	email   string
}

// Name returns the name of the Violation.
func (*StaleCommitViolation) Name() string {
	return "StaleCommitViolation"
}

// Message implements Violation.
func (sc *StaleCommitViolation) Message() string {
	return sc.message
}

// FileLocation implements Violation.
func (*StaleCommitViolation) FileLocation() (string, error) {
	return "", ErrViolationMethodNotExist
}

// LineLocation implements Violation.
func (*StaleCommitViolation) LineLocation() (int, error) {
	return 0, ErrViolationMethodNotExist
}

// Suggestion implements Violation.
func (*StaleCommitViolation) Suggestion() (string, error) {
	return "", ErrViolationMethodNotExist
}

// Author implements Violation.
func (p *StaleCommitViolation) Author() (*remote.Author, error) {
	return nil, ErrViolationMethodNotExist
}

// Severity implements Violation.
func (p *StaleCommitViolation) Severity() Severity {
	return Violated
}

// Email implements Violation.
func (sc *StaleCommitViolation) Email() string {
	return sc.email
}
