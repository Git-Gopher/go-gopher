package violation

import (
	"fmt"
	"time"

	"github.com/Git-Gopher/go-gopher/model/github"
)

func NewStaleBranchViolation(branch string, duration time.Duration) *StaleBranchViolation {
	stale := &StaleBranchViolation{
		display:  nil,
		branch:   branch,
		duration: duration,
	}
	stale.display = &display{stale}

	return stale
}

// Example violation.
type StaleBranchViolation struct {
	*display
	branch   string
	duration time.Duration
}

// Name returns the name of the Violation.
func (*StaleBranchViolation) Name() string {
	return "StaleBranchViolation"
}

// Message implements Violation.
func (sbv *StaleBranchViolation) Message() string {
	return fmt.Sprintf("Branch \"%s\" is stale due it not being committed to for over %d months", sbv.branch, sbv.duration)
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
	return fmt.Sprintf("Consider deleting the branch, \"%s\" if it is unused delete it, or continue to work use it by merging the primary branch into it", sbv.branch), nil
}

// Author implements Violation.
func (*StaleBranchViolation) Author() (*github.Author, error) {
	return nil, ErrViolationMethodNotExist
}

// Severity implements Violation.
func (p *StaleBranchViolation) Severity() Severity {
	return Violated
}
