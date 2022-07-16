package violation

import (
	"time"

	"github.com/Git-Gopher/go-gopher/model/remote"
)

func NewCommonViolation(message string, author *remote.Author, time time.Time) *CommonViolation {
	common := &CommonViolation{
		violation: violation{
			name:     "CommonViolation",
			email:    author.Email,
			time:     time,
			severity: Suggestion,
		},
		message: message,
		author:  author,
	}
	common.display = &display{common}

	return common
}

// Example violation.
type CommonViolation struct {
	violation
	*display
	message string
	author  *remote.Author
}

// Message implements Violation.
func (cv *CommonViolation) Message() string {
	return cv.message
}

// Author implements Violation.
func (cv *CommonViolation) Author() (*remote.Author, error) {
	return cv.author, nil
}
