package violation

import (
	"fmt"
	"strings"

	"github.com/Git-Gopher/go-gopher/model/github"
)

func NewDescriptiveCommitViolation(
	message string,
	email string,
) *DescriptiveCommitViolation {
	violation := &DescriptiveCommitViolation{
		display: nil,
		message: message,
		email:   email,
	}
	violation.display = &display{violation}

	return violation
}

// from feature branches.
type DescriptiveCommitViolation struct {
	*display
	message string
	email   string
}

// Name implements Violation.
func (dvc *DescriptiveCommitViolation) Name() string {
	return "DescriptiveCommitViolation"
}

// Message implements Violation.
func (dvc *DescriptiveCommitViolation) Message() string {
	message := strings.ReplaceAll(dvc.message, "\n", " ")

	return fmt.Sprintf(`The commit message \"%s\" may not be 
		descriptive enough`, message)
}

// Suggestion implements Violation.
func (dvc *DescriptiveCommitViolation) Suggestion() (string, error) {
	return "Try to add more detail to commit messages that relate to the content of a commit", nil
}

// Author implements Violation.
func (dvc *DescriptiveCommitViolation) Author() (*github.Author, error) {
	return &github.Author{
		Email:     dvc.email,
		Login:     "",
		AvatarUrl: "",
	}, nil
}

// FileLocation implements Violation.
func (dvc *DescriptiveCommitViolation) FileLocation() (string, error) {
	return "", ErrViolationMethodNotExist
}

// LineLocation implements Violation.
func (dvc *DescriptiveCommitViolation) LineLocation() (int, error) {
	return 0, ErrViolationMethodNotExist
}

// Severity implements Violation.
func (dvc *DescriptiveCommitViolation) Severity() Severity {
	return Suggestion
}

// Email implements Violation.
func (dvc *DescriptiveCommitViolation) Email() string {
	return dvc.email
}
