package violation

import (
	"fmt"
	"strings"

	"github.com/Git-Gopher/go-gopher/model/remote"
	"github.com/Git-Gopher/go-gopher/utils"
)

func NewForcePushViolation(
	lostCommits []utils.Commit,
	email string,
) *ForcePushViolation {
	violation := &ForcePushViolation{
		display:     nil,
		lostCommits: lostCommits,
		email:       email,
	}
	violation.display = &display{violation}

	return violation
}

// Force push violation occurs whenever a branch is force pushed to, losing a series of commits
// from feature branches.
type ForcePushViolation struct {
	*display
	lostCommits []utils.Commit
	email       string
}

// FileLocation implements Violation.
func (*ForcePushViolation) FileLocation() (string, error) {
	return "", ErrViolationMethodNotExist
}

// LineLocation implements Violation.
func (*ForcePushViolation) LineLocation() (int, error) {
	return 0, ErrViolationMethodNotExist
}

// Message implements Violation.
func (f *ForcePushViolation) Message() string {
	format := "The following commits have been lost:\n%s"
	commits := make([]string, len(f.lostCommits))
	for i, commit := range f.lostCommits {
		commits[i] = commit.Link()
	}

	return fmt.Sprintf(format, strings.Join(commits, ",\n"))
}

// Name implements Violation.
func (*ForcePushViolation) Name() string {
	return "ForcePushViolation"
}

// Suggestion implements Violation.
func (f *ForcePushViolation) Suggestion() (string, error) {
	format := "Restore the following commits to restore the work lost on the branch:\n%s"
	commits := make([]string, len(f.lostCommits))
	for i, commit := range f.lostCommits {
		commits[i] = commit.Link()
	}

	return fmt.Sprintf(format, strings.Join(commits, ",\n")), nil
}

// Author implements Violation.
func (f *ForcePushViolation) Author() (*remote.Author, error) {
	return nil, ErrViolationMethodNotExist
}

// Severity implements Violation.
func (f *ForcePushViolation) Severity() Severity {
	return Violated
}

// Email implements Violation.
func (f *ForcePushViolation) Email() string {
	return f.email
}
