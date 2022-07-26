package violation

import (
	"fmt"
	"time"
)

func NewBinaryViolation(
	filename string,
	email string,
	time time.Time,
) *BinaryViolation {
	violation := &BinaryViolation{
		violation: violation{
			name:     "BinaryViolation",
			email:    email,
			time:     time,
			severity: Suggestion,
		},
		filename: filename,
	}
	violation.display = &display{violation}

	return violation
}

// BinaryViolation is violation when a branch name is inconsistent with others.
// from feature branches.
type BinaryViolation struct {
	violation
	*display
	filename string
}

// Message implements Violation.
func (bv *BinaryViolation) Message() string {
	format := "A binary file \"%s\" has been committed to the project"

	return fmt.Sprintf(format, bv.name)
}

// Suggestion implements Violation.
func (bv *BinaryViolation) Suggestion() (string, error) {
	return `Git projects should aim to not include binary files.
	  Binary files should be added to the project .gitignore file and the 
	  file removed from the working tree using \"git rm --cached <file>\"`, nil
}
