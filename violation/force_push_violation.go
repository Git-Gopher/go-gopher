package violation

import "fmt"

func NewForcePushViolation(
	lostCommits []string,
) *ForcePushViolation {
	violation := &ForcePushViolation{
		display:     nil,
		lostCommits: lostCommits,
	}
	violation.display = &display{violation}

	return violation
}

// Force push violation occurs whenever a branch is force pushed to, losing a series of commits
// from feature branches.
type ForcePushViolation struct {
	*display
	lostCommits []string
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
	format := "The following commits have been lost: %v"

	return fmt.Sprintf(format, f.lostCommits)
}

// Name implements Violation.
func (*ForcePushViolation) Name() string {
	return "ForcePushViolation"
}

// Suggestion implements Violation.
func (f *ForcePushViolation) Suggestion() (string, error) {
	return fmt.Sprintf("Restore the following commits to restore the work lost on the branch: \"%v\" ", f.lostCommits), nil
}
