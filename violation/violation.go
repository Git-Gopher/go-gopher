package violation

import (
	"errors"
	"fmt"
	"time"

	"github.com/Git-Gopher/go-gopher/model/remote"
)

type Severity int

const (
	Violated Severity = iota
	Suggestion
)

var ErrViolationMethodNotExist = errors.New("Violation method not exist")

type Violation interface {
	Name() string       // required: Internal name of the violation.
	Message() string    // required: Warning message.
	Display() string    // required: Formal display line of the violation.
	Email() string      // required: Email address of the violator.
	Time() time.Time    // required: Time of the violation.
	Severity() Severity // required: Severity of the violation.

	Author() (*remote.Author, error) // optional: GitHub author which caused the violation.
	FileLocation() (string, error)   // optional: File location of the violation.
	LineLocation() (int, error)      // optional: Line location of violation.
	Suggestion() (string, error)     // optional: Suggestion to fix the violation.
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

// violation - common base struct for all violations.
type violation struct {
	name     string
	email    string
	time     time.Time
	severity Severity
}

func (v *violation) Name() string {
	return v.name
}

func (v *violation) Email() string {
	return v.email
}

func (v *violation) Time() time.Time {
	return v.time
}

func (v *violation) Severity() Severity {
	return v.severity
}

func (v *violation) Author() (*remote.Author, error) {
	return nil, ErrViolationMethodNotExist
}

func (v *violation) FileLocation() (string, error) {
	return "", ErrViolationMethodNotExist
}

func (v *violation) LineLocation() (int, error) {
	return 0, ErrViolationMethodNotExist
}

func (v *violation) Suggestion() (string, error) {
	return "", ErrViolationMethodNotExist
}
