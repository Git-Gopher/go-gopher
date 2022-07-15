package detector

import (
	"strings"

	"github.com/Git-Gopher/go-gopher/markup"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/violation"
)

type CommitDetect func(c *common, commit *local.Commit) (bool, violation.Violation, error)

type CommitDetector struct {
	violated   int
	found      int
	total      int
	violations []violation.Violation

	detect CommitDetect
}

func (cd *CommitDetector) Run(model *enriched.EnrichedModel) error {
	if model == nil {
		return nil
	}

	// Struct should be reset before each run, incase we are running it with a different model.
	cd.violated = 0
	cd.found = 0
	cd.total = 0
	cd.violations = make([]violation.Violation, 0)

	c := common{owner: model.Owner, repo: model.Name}

	for _, co := range model.Commits {
		co := co
		detected, violation, err := cd.detect(&c, &co)
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
func BranchCommitDetect() CommitDetect {
	return func(c *common, commit *local.Commit) (bool, violation.Violation, error) {
		if len(commit.ParentHashes) >= 2 {
			return true, nil, nil
		}

		return false, nil, nil
	}
}

// XXX: Very very lazy. I am a true software engineer.
func DiffMatchesMessageDetect() CommitDetect {
	return func(c *common, commit *local.Commit) (bool, violation.Violation, error) {
		words := strings.Split(commit.Message, " ")
		for _, diff := range commit.DiffToParents {
			for _, word := range words {
				all := diff.Addition + diff.Deletion + diff.Equal + diff.Name
				if strings.Contains(all, word) {
					return true, nil, nil
				}
			}
		}

		return false, violation.NewDescriptiveCommitViolation(markup.Commit{
			Hash: commit.Hash.String(),
			GitHubLink: markup.GitHubLink{
				Owner: c.owner,
				Repo:  c.repo,
			},
		}, commit.Message, commit.Author.Email), nil
	}
}

// Check if commit is less than 3 words.
func ShortCommitMessageDetect() CommitDetect {
	return func(c *common, commit *local.Commit) (bool, violation.Violation, error) {
		exclusions := []string{
			"first commit",
			"initial commit",
		}
		for _, exclusion := range exclusions {
			if strings.ToLower(commit.Message) == exclusion {
				return true, nil, nil
			}
		}

		words := strings.Split(commit.Message, " ")
		if len(words) < 3 {
			return false, violation.NewShortCommitViolation(
				markup.Commit{
					Hash: commit.Hash.String(),
					GitHubLink: markup.GitHubLink{
						Owner: c.owner,
						Repo:  c.repo,
					},
				},
				commit.Message,
				commit.Author.Email,
			), nil
		}

		return true, nil, nil
	}
}
