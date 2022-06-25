package violation

import "github.com/Git-Gopher/go-gopher/model/github"

func NewCommonViolation(message string, author *github.Author) *CommonViolation {
	common := &CommonViolation{display: nil, message: message, author: author}
	common.display = &display{common}

	return common
}

// Example violation.
type CommonViolation struct {
	*display
	message string
	author  *github.Author
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

// Suggestion implements Violation.
func (cv *CommonViolation) Author() (*github.Author, error) {
	return cv.author, nil
}
