package violation

import (
	"errors"
	"fmt"
	"time"

	"github.com/Git-Gopher/go-gopher/markup"
	"github.com/Git-Gopher/go-gopher/utils"
)

type Severity int

const (
	Violated Severity = iota
	Suggestion
)

var (
	ErrViolationMethodNotExist = errors.New("Violation method not exist")
	ErrCreatedTimePullRequest  = errors.New("no created time for pull request")
	ErrClosedTimePullRequest   = errors.New("no closed time for pull request")
)

type Violation interface {
	Name() string                 // required: Internal name of the violation.
	Message() string              // required: Warning message.
	Display(utils.Authors) string // required: Formal display line of the violation.
	Time() time.Time              // required: Time of the violation.
	Severity() Severity           // required: Severity of the violation.
	Current() bool                // required: Is the violation related to the current reporting

	Email() (string, error)        // optional Email address of the violator.
	Login() (string, error)        // optional username of the violator.
	FileLocation() (string, error) // optional: File location of the violation.
	LineLocation() (int, error)    // optional: Line location of violation.
	Suggestion() (string, error)   // optional: Suggestion to fix the violation.
}

// display struct implements the Display method part of Violation using Violation.
type display struct {
	v Violation
}

// Display implements Violation.
func (d *display) Display(authors utils.Authors) string {
	// Get the author of the violation.
	authorLink := "unknown"
	if login, err := d.v.Login(); err == nil {
		authorLink = markup.Author(login).Markdown()
	} else if email, err := d.v.Email(); err == nil {
		login, err := authors.Find(email)
		if login != nil && err == nil {
			authorLink = markup.Author(*login).Markdown()
		}
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
	login    string
	time     time.Time
	severity Severity
	current  bool
}

func (v *violation) Name() string {
	return v.name
}

func (v *violation) Email() (string, error) {
	if v.email == "" {
		return "", ErrViolationMethodNotExist
	}

	return v.email, nil
}

func (v *violation) Login() (string, error) {
	if v.login == "" {
		return "", ErrViolationMethodNotExist
	}

	return v.login, nil
}

func (v *violation) Time() time.Time {
	return v.time
}

func (v *violation) Severity() Severity {
	return v.severity
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

func FilterByLogin(violations []Violation, users utils.Authors, filter []string) []Violation {
	if len(filter) == 0 || len(violations) == 0 {
		return violations
	}
	loginMap := make(map[string]struct{})
	for _, l := range filter {
		loginMap[l] = struct{}{}
	}

	filtered := []Violation{}
	for _, v := range violations {
		login := ""
		if l, err := v.Login(); err == nil {
			login = l
		} else if email, err := v.Email(); users != nil && err == nil {
			if l, err := users.Find(email); err == nil {
				login = *l
			}
		} else {
			// Violation does not have author.
			filtered = append(filtered, v)

			continue
		}

		if _, ok := loginMap[login]; !ok {
			filtered = append(filtered, v)
		}
	}

	return filtered
}
