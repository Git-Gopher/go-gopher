package violation

import (
	"fmt"
	"time"

	"github.com/Git-Gopher/go-gopher/markup"
)

func NewExtremeDiffDistanceViolation(commit markup.Commit, email string, time time.Time) *DiffDistanceViolation {
	common := &DiffDistanceViolation{
		violation: violation{
			name:     "ExtremeDiffDistanceViolation",
			email:    email,
			time:     time,
			severity: Violated,
		},
		commit: commit,
	}
	common.display = &display{common}

	return common
}

func NewMildDiffDistanceViolation(commit markup.Commit, email string, time time.Time) *DiffDistanceViolation {
	common := &DiffDistanceViolation{
		violation: violation{
			name:     "MildDiffDistanceViolation",
			email:    email,
			time:     time,
			severity: Violated,
		},
		commit: commit,
	}
	common.display = &display{common}

	return common
}

// Example violation.
type DiffDistanceViolation struct {
	violation
	*display
	commit markup.Commit
}

// Message implements Violation.
func (sc *DiffDistanceViolation) Message() string {
	if sc.severity == Violated {
		return fmt.Sprintf("Diff distance extreme violation on %s", sc.commit.Markdown())
	} else {
		return fmt.Sprintf("Diff distance mild violation on %s", sc.commit.Markdown())
	}
}
