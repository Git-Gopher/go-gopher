package violation

import (
	"errors"
	"fmt"
)

var ErrViolationMethodNotExist = errors.New("Violation method not exist")

type Violation interface {
	Name() string                // required: Internal name of the violation.
	Message() string             // required: Warning message.
	Suggestion() (string, error) // required: Suggests a remedy for the violation.
	Display() string             // required: Formal display line of the violation
	FileLocation() (string, error)
	LineLocation() (int, error)
}

// display struct implements the Display method part of Violation using Violation.
type display struct {
	v Violation
}

// Display implements Violation.
func (d *display) Display() string {
	var format string = "%s: %s"
	suggestion, err := d.v.Suggestion()
	if err != nil {
		return fmt.Sprintf(format, d.v.Name(), d.v.Message())
	}
	format += "\n\t%s"

	return fmt.Sprintf(format, d.v.Name(), d.v.Message(), suggestion)
}
