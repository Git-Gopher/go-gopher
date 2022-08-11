package violation

import (
	"fmt"
	"time"

	"github.com/Git-Gopher/go-gopher/markup"
)

func NewStaleBranchViolation(branch markup.Branch, duration time.Duration, email string) *StaleBranchViolation {
	stale := &StaleBranchViolation{
		violation: violation{
			name:     "StaleBranchViolation",
			email:    email,
			time:     time.Now(),
			severity: Violated,
		},
		branch:   branch,
		duration: duration,
	}
	stale.display = &display{stale}

	return stale
}

// Example violation.
type StaleBranchViolation struct {
	violation
	*display
	branch   markup.Branch
	duration time.Duration
}

// Message implements Violation.
func (sbv *StaleBranchViolation) Message() string {
	return fmt.Sprintf(
		"Branch %s is stale due it not being committed to for over %d weeks",
		sbv.branch.Markdown(),
		sbv.duration,
	)
}

// Suggestion implements Violation.
func (sbv *StaleBranchViolation) Suggestion() (string, error) {
	return fmt.Sprintf("Consider deleting the branch, \"%s\" if it is unused delete it,"+
		" or continue to work use it by merging the primary branch into it", sbv.branch), nil
}
