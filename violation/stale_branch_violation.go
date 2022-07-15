package violation

import (
	"fmt"
	"time"

	"github.com/Git-Gopher/go-gopher/markup"
	"github.com/Git-Gopher/go-gopher/model/remote"
)

func NewStaleBranchViolation(branch markup.Branch, duration time.Duration, email string) *StaleBranchViolation {
	stale := &StaleBranchViolation{
		display:  nil,
		branch:   branch,
		duration: duration,
		email:    email,
	}
	stale.display = &display{stale}

	return stale
}

// Example violation.
type StaleBranchViolation struct {
	*display
	branch   markup.Branch
	duration time.Duration
	email    string
}

// Name returns the name of the Violation.
func (*StaleBranchViolation) Name() string {
	return "StaleBranchViolation"
}

// Message implements Violation.
func (sbv *StaleBranchViolation) Message() string {
	return fmt.Sprintf(
		"Branch \"%s\" is stale due it not being committed to for over %d months",
		sbv.branch.Link(),
		sbv.duration,
	)
}

// FileLocation implements Violation.
func (*StaleBranchViolation) FileLocation() (string, error) {
	return "", ErrViolationMethodNotExist
}

// LineLocation implements Violation.
func (*StaleBranchViolation) LineLocation() (int, error) {
	return 0, ErrViolationMethodNotExist
}

// Suggestion implements Violation.
func (sbv *StaleBranchViolation) Suggestion() (string, error) {
	return fmt.Sprintf(`Consider deleting the branch, \"%s\" if it 
		is unused delete it, or continue to work use it by merging the primary branch into it`, sbv.branch), nil
}

// Author implements Violation.
func (*StaleBranchViolation) Author() (*remote.Author, error) {
	return nil, ErrViolationMethodNotExist
}

// Severity implements Violation.
func (sbv *StaleBranchViolation) Severity() Severity {
	return Violated
}

// Email implements Violation.
func (sbv *StaleBranchViolation) Email() string {
	return sbv.email
}
