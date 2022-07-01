// nolint:dupl
package violation

import (
	"fmt"
	"strings"

	"github.com/Git-Gopher/go-gopher/model/github"
)

func NewShortCommitViolation(
	message string,
	author string,
) *ShortCommitViolation {
	violation := &ShortCommitViolation{
		display: nil,
		message: message,
		author:  author,
	}
	violation.display = &display{violation}

	return violation
}

// from feature branches.
type ShortCommitViolation struct {
	*display
	message string
	author  string
}

// Name implements Violation.
func (dvc *ShortCommitViolation) Name() string {
	return "ShortCommitViolation"
}

// Message implements Violation.
func (dvc *ShortCommitViolation) Message() string {
	message := strings.ReplaceAll(dvc.message, "\n", " ")

	return fmt.Sprintf(`Commit message \"%s\" is too short`, message)
}

// Suggestion implements Violation.
func (dvc *ShortCommitViolation) Suggestion() (string, error) {
	return "Commit message is too short", nil
}

// Author implements Violation.
func (dvc *ShortCommitViolation) Author() (*github.Author, error) {
	return &github.Author{
		Login:     dvc.author,
		AvatarUrl: "",
	}, nil
}

// FileLocation implements Violation.
func (dvc *ShortCommitViolation) FileLocation() (string, error) {
	return "", ErrViolationMethodNotExist
}

// LineLocation implements Violation.
func (dvc *ShortCommitViolation) LineLocation() (int, error) {
	return 0, ErrViolationMethodNotExist
}

// Severity implements Violation.
func (p *ShortCommitViolation) Severity() Severity {
	return Suggestion
}
