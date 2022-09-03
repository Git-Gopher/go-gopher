package violation

import (
	"errors"
	"fmt"
	"time"

	"github.com/Git-Gopher/go-gopher/markup"
	"github.com/Git-Gopher/go-gopher/model/remote"
	"github.com/Git-Gopher/go-gopher/utils"
)

type Severity int

const (
	Violated Severity = iota
	Suggestion
)

var ErrViolationMethodNotExist = errors.New("Violation method not exist")

type Violation interface {
	Name() string                 // required: Internal name of the violation.
	Message() string              // required: Warning message.
	Display(utils.Authors) string // required: Formal display line of the violation.
	Email() string                // required: Email address of the violator.
	Time() time.Time              // required: Time of the violation.
	Severity() Severity           // required: Severity of the violation.
	Current() bool                // required: Is the violation related to the current reporting

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
func (d *display) Display(authors utils.Authors) string {
	// Get the author of the violation.
	authorLink := "unknown"
	author, _ := authors.Find(d.v.Email())
	if author != nil {
		authorLink = markup.Author(*author).Markdown()
	}

	suggestion, err := d.v.Suggestion()
	if err != nil {
		// If the suggestion is not available.
		return fmt.Sprintf(
			"%s: %s - %s %s\n",
			d.v.Name(),
			d.v.Message(),
			authorLink,
			d.v.Time().Format(time.UnixDate),
		)
	}

	return fmt.Sprintf(
		"%s: %s - %s %s \n\t%s\n",
		d.v.Name(),
		d.v.Message(),
		authorLink,
		d.v.Time().Format(time.UnixDate),
		suggestion,
	)
}

// violation - common base struct for all violations.
type violation struct {
	name     string
	email    string
	time     time.Time
	severity Severity
	current  bool
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

func (v *violation) Current() bool {
	return v.current
}

func FilterByLogin(violations []Violation, login []string) []Violation {
	if len(login) == 0 || len(violations) == 0 {
		return violations
	}

	loginMap := make(map[string]struct{})
	for _, l := range login {
		loginMap[l] = struct{}{}
	}

	filtered := []Violation{}
	for _, v := range violations {
		author, err := v.Author()
		if err != nil {
			filtered = append(filtered, v)

			continue
		}

		if _, ok := loginMap[author.Login]; !ok {
			filtered = append(filtered, v)
		}
	}

	return filtered
}
