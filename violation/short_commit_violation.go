package violation

import (
	"fmt"
	"strings"

	"github.com/Git-Gopher/go-gopher/model/remote"
	"github.com/Git-Gopher/go-gopher/utils"
)

func NewShortCommitViolation(
	commit utils.Commit,
	message string,
	email string,
) *ShortCommitViolation {
	violation := &ShortCommitViolation{
		display: nil,
		commit:  commit,
		message: message,
		email:   email,
	}
	violation.display = &display{violation}

	return violation
}

// from feature branches.
type ShortCommitViolation struct {
	*display
	commit  utils.Commit
	message string
	email   string
}

// Name implements Violation.
func (sc *ShortCommitViolation) Name() string {
	return "ShortCommitViolation"
}

// Message implements Violation.
func (sc *ShortCommitViolation) Message() string {
	message := strings.ReplaceAll(sc.message, "\n", " ")

	return fmt.Sprintf(`Commit message \"%s\" on \"%s\" is too short`, message, sc.commit.Link())
}

// Suggestion implements Violation.
func (sc *ShortCommitViolation) Suggestion() (string, error) {
	return "Commit message is too short", nil
}

// Author implements Violation.
func (sc *ShortCommitViolation) Author() (*remote.Author, error) {
	return nil, ErrViolationMethodNotExist
}

// FileLocation implements Violation.
func (sc *ShortCommitViolation) FileLocation() (string, error) {
	return "", ErrViolationMethodNotExist
}

// LineLocation implements Violation.
func (sc *ShortCommitViolation) LineLocation() (int, error) {
	return 0, ErrViolationMethodNotExist
}

// Severity implements Violation.
func (sc *ShortCommitViolation) Severity() Severity {
	return Suggestion
}

// Email implements Violation.
func (sc *ShortCommitViolation) Email() string {
	return sc.email
}
