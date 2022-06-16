package detector

import (
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/violation"
)

// Example detector to check if the commit is greater than 10 characters.
func NewDecriptiveCommitMessageDetect() CommitDetect {
	return func(commit *local.Commit) (bool, violation.Violation, error) {
		return false, nil, ErrNotImplemented
	}
}
