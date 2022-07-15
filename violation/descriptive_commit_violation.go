package violation

import (
	"fmt"
	"strings"
	"time"

	"github.com/Git-Gopher/go-gopher/markup"
)

func NewDescriptiveCommitViolation(
	commit markup.Commit,
	message string,
	email string,
	time time.Time,
) *DescriptiveCommitViolation {
	violation := &DescriptiveCommitViolation{
		violation: violation{
			name:     "DescriptiveCommitViolation",
			email:    email,
			time:     time,
			severity: Suggestion,
		},
		commit:  commit,
		message: message,
		email:   email,
	}
	violation.display = &display{violation}

	return violation
}

// from feature branches.
type DescriptiveCommitViolation struct {
	violation
	*display
	commit  markup.Commit
	message string
	email   string
}

// Message implements Violation.
func (dvc *DescriptiveCommitViolation) Message() string {
	message := strings.ReplaceAll(dvc.message, "\n", " ")

	return fmt.Sprintf(`The commit message \"%s\" on \"%s\" may not be 
		descriptive enough`, message, dvc.commit.Link())
}

// Suggestion implements Violation.
func (dvc *DescriptiveCommitViolation) Suggestion() (string, error) {
	return "Try to add more detail to commit messages that relate to the content of a commit", nil
}
