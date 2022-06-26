package detector

import (
	"encoding/hex"
	"strings"

	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/violation"
)

type CommitDetect func(commit *local.Commit) (bool, violation.Violation, error)

type CommitDetector struct {
	violated   int
	found      int
	total      int
	violations []violation.Violation

	detect CommitDetect
}

func (cd *CommitDetector) Run(model *enriched.EnrichedModel) error {
	// Struct should be reset before each run, incase we are running it with a different model.
	cd.violated = 0
	cd.found = 0
	cd.total = 0
	cd.violations = make([]violation.Violation, 0)

	for _, c := range model.Commits {
		c := c
		detected, violation, err := cd.detect(&c)
		cd.total++
		if err != nil {
			return err
		}
		if detected {
			cd.found++
		}
		if violation != nil {
			cd.violations = append(cd.violations, violation)
		}
	}

	return nil
}

func (cd *CommitDetector) Result() (int, int, int, []violation.Violation) {
	return cd.violated, cd.found, cd.total, cd.violations
}

func NewCommitDetector(detect CommitDetect) *CommitDetector {
	return &CommitDetector{
		violated:   0,
		found:      0,
		total:      0,
		violations: make([]violation.Violation, 0),
		detect:     detect,
	}
}

// All commits on the main branch for github flow should be merged in,
// meaning that they have two parents(the main branch and the feature branch).
func TwoParentsCommitDetect() CommitDetect {
	return func(commit *local.Commit) (bool, violation.Violation, error) {
		if len(commit.ParentHashes) >= 2 {
			return true, nil, nil
		}

		commitHash := hex.EncodeToString(commit.Hash.ToByte())

		parentHashes := make([]string, len(commit.ParentHashes))
		for i, p := range commit.ParentHashes {
			parentHashes[i] = hex.EncodeToString(p.ToByte())
		}

		// TODO: Don't use hardcoded primary branch
		return false, violation.NewPrimaryBranchDirectCommitViolation("main", commitHash, parentHashes), nil
	}
}

// XXX: Very very lazy. I am a true software engineer.
func DiffMatchesMessageDetect() CommitDetect {
	return func(commit *local.Commit) (bool, violation.Violation, error) {
		words := strings.Split(commit.Message, " ")
		for _, diff := range commit.DiffToParents {
			for _, word := range words {
				all := diff.Addition + diff.Deletion + diff.Equal + diff.Name
				if strings.Contains(all, word) {
					return true, nil, nil
				}
			}
		}

		return false, violation.NewDescriptiveCommitViolation(commit.Message, commit.Author.Name), nil
	}
}

// Example detector to check if the commit is greater than 10 characters.
func NewDecriptiveCommitMessageDetect() CommitDetect {
	return func(commit *local.Commit) (bool, violation.Violation, error) {
		return false, nil, ErrNotImplemented
	}
}
