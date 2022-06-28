package violation

import (
	"errors"
	"fmt"

	"github.com/Git-Gopher/go-gopher/model/github"
)

type Severity int

const (
	Violated Severity = iota
	Suggestion
)

var ErrViolationMethodNotExist = errors.New("Violation method not exist")

type Violation interface {
	Name() string                // required: Internal name of the violation.
	Message() string             // required: Warning message.
	Suggestion() (string, error) // required: Suggests a remedy for the violation.
	Display() string             // required: Formal display line of the violation
	// XXX: Use enriched model instead.
	Author() (*github.Author, error) // optional: User which caused the violation the
	FileLocation() (string, error)
	LineLocation() (int, error)
	Severity() Severity
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
	format += "\n\t%s\n"

	return fmt.Sprintf(format, d.v.Name(), d.v.Message(), suggestion)
}
