package violation

import (
	"fmt"
	"time"

	"github.com/Git-Gopher/go-gopher/markup"
)

func NewBinaryViolation(
	file markup.File,
	email string,
	time time.Time,
) *BinaryViolation {
	violation := &BinaryViolation{
		violation: violation{
			name:     "BinaryViolation",
			email:    email,
			time:     time,
			severity: Violated,
		},
		file: file,
	}
	violation.display = &display{violation}

	return violation
}

// BinaryViolation is violation when a binary has been committed to the repository.
type BinaryViolation struct {
	violation
	*display
	file markup.File
}

// Message implements Violation.
func (bv *BinaryViolation) Message() string {
	format := "A binary file %s has been committed to the project"

	return fmt.Sprintf(format, bv.file.Markdown())
}

// Suggestion implements Violation.
func (bv *BinaryViolation) Suggestion() (string, error) {
	return `Git projects should aim to not include binary files.
	  Binary files should be added to the project .gitignore file and the 
	  file removed from the working tree using \"git rm --cached <file>\"`, nil
}
