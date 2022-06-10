package detector

import (
	"github.com/Git-Gopher/go-gopher/model/github"
	"github.com/Git-Gopher/go-gopher/model/local"
)

type CommitDetect func(commit *local.Commit) (bool, error)

type CommitDetector struct {
	violated int
	found    int
	total    int

	detect CommitDetect
}

// TODO: We should change this to the enriched model
func (cd *CommitDetector) Run(model *local.GitModel) error {
	for _, c := range model.Commits {
		c := c
		detected, err := cd.detect(&c)
		cd.total++
		if err != nil {
			return err
		}
		if detected {
			cd.found++
		}
	}

	return nil
}

func (cd *CommitDetector) Run2(model *github.GithubModel) error {
	return nil
}

func (cd *CommitDetector) Result() (int, int, int) {
	return cd.violated, cd.found, cd.total
}

func NewCommitDetector(detect CommitDetect) *CommitDetector {
	return &CommitDetector{
		violated: 0,
		found:    0,
		total:    0,
		detect:   detect,
	}
}

func NewLineLengthCommitDetect() CommitDetect {
	return func(commit *local.Commit) (bool, error) {
		if len(commit.Message) > 10 {
			return false, nil
		}

		return true, nil
	}
}

// All commits on the main branch for github flow should be merged in,
// meaning that they have two parents(the main branch and the feature branch).
func TwoParentsCommitDetect() CommitDetect {
	return func(commit *local.Commit) (bool, error) {
		if len(commit.ParentHashes) >= 2 {
			return true, nil
		}

		return false, nil
	}
}
