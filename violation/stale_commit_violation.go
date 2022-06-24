package violation

import "github.com/Git-Gopher/go-gopher/model/github"

func NewStaleCommitViolation(message string) *StaleCommitViolation {
	common := &StaleCommitViolation{display: nil, message: message}
	common.display = &display{common}

	return common
}

// Example violation.
type StaleCommitViolation struct {
	*display
	message string
}

// Name returns the name of the Violation.
func (*StaleCommitViolation) Name() string {
	return "StaleCommitViolation"
}

// Message implements Violation.
func (cv *StaleCommitViolation) Message() string {
	return cv.message
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
func (p *StaleCommitViolation) Author() (*github.Author, error) {
	return nil, ErrViolationMethodNotExist
}
