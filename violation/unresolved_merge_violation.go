package violation

import (
	"fmt"
	"strings"

	"github.com/Git-Gopher/go-gopher/model/github"
)

func NewUnresolvedMergeViolation(
	commit string,
	filepath string,
	line int,
	email string,
) *UnresolvedMergeViolation {
	violation := &UnresolvedMergeViolation{
		display:  nil,
		filepath: filepath,
		line:     line,
		email:    email,
	}
	violation.display = &display{violation}

	return violation
}

// Violation when unresolved merge conflicts are merged violation.
type UnresolvedMergeViolation struct {
	*display
	commit   string
	filepath string
	line     int
	email    string
}

// Name implements Violation.
func (um *UnresolvedMergeViolation) Name() string {
	return "UnresolvedMergeViolation"
}

// Message implements Violation.
func (um *UnresolvedMergeViolation) Message() string {
	message := strings.ReplaceAll(um.message, "\n", " ")

	return fmt.Sprintf(`The commit message \"%s\" may not be 
		descriptive enough`, message)
}

// Suggestion implements Violation.
func (um *UnresolvedMergeViolation) Suggestion() (string, error) {
	return "Try to add more detail to commit messages that relate to the content of a commit", nil
}

// Author implements Violation.
func (um *UnresolvedMergeViolation) Author() (*github.Author, error) {
	return &github.Author{
		Email:     um.email,
		Login:     "",
		AvatarUrl: "",
	}, nil
}

// FileLocation implements Violation.
func (um *UnresolvedMergeViolation) FileLocation() (string, error) {
	return "", ErrViolationMethodNotExist
}

// LineLocation implements Violation.
func (um *UnresolvedMergeViolation) LineLocation() (int, error) {
	return 0, ErrViolationMethodNotExist
}

// Severity implements Violation.
func (um *UnresolvedMergeViolation) Severity() Severity {
	return Suggestion
}

// Email implements Violation.
func (um *UnresolvedMergeViolation) Email() string {
	return um.email
}
