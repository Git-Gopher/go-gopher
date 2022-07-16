package violation

import (
	"fmt"
	"strings"
	"time"

	"github.com/Git-Gopher/go-gopher/markup"
)

func NewShortCommitViolation(
	commit markup.Commit,
	message string,
	email string,
	time time.Time,
) *ShortCommitViolation {
	violation := &ShortCommitViolation{
		violation: violation{
			name:     "ShortCommitViolation",
			email:    email,
			time:     time,
			severity: Violated,
		},
		commit:  commit,
		message: message,
	}
	violation.display = &display{violation}

	return violation
}

// from feature branches.
type ShortCommitViolation struct {
	violation
	*display
	commit  markup.Commit
	message string
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
