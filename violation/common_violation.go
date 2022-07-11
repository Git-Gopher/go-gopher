package violation

import (
	"github.com/Git-Gopher/go-gopher/model/remote"
)

func NewCommonViolation(message string, author *remote.Author) *CommonViolation {
	common := &CommonViolation{display: nil, message: message, author: author}
	common.display = &display{common}

	return common
}

// Example violation.
type CommonViolation struct {
	*display
	message string
	author  *remote.Author
}

// Name returns the name of the Violation.
func (*CommonViolation) Name() string {
	return "CommonViolation"
}

// Message implements Violation.
func (cv *CommonViolation) Message() string {
	return cv.message
}

// FileLocation implements Violation.
func (*CommonViolation) FileLocation() (string, error) {
	return "", ErrViolationMethodNotExist
}

// LineLocation implements Violation.
func (*CommonViolation) LineLocation() (int, error) {
	return 0, ErrViolationMethodNotExist
}

// Suggestion implements Violation.
func (*CommonViolation) Suggestion() (string, error) {
	return "", ErrViolationMethodNotExist
}

// Author implements Violation.
func (cv *CommonViolation) Author() (*remote.Author, error) {
	return cv.author, nil
}

// Severity implements Violation.
func (*CommonViolation) Severity() Severity {
	return Suggestion
}

// Email implements Violation.
func (cv *CommonViolation) Email() string {
	return cv.author.Email
}
