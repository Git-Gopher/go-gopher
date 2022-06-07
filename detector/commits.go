package detector

import (
	"github.com/Git-Gopher/go-gopher/model"
)

type CommitDetect func(commit *model.Commit) (bool, error)

type CommitDetector struct {
	violated int
	found    int
	total    int

	detect CommitDetect
}

func (cd *CommitDetector) Run(model *model.GitModel) error {
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

func (c *CommitDetector) Result() (int, int, int) {
	return c.violated, c.found, c.total
}

func NewCommitDetector(detect CommitDetect) *CommitDetector {
	return &CommitDetector{
		total:  0,
		found:  0,
		detect: detect,
	}
}

func NewLineLengthCommitDetect() CommitDetect {
	return func(commit *model.Commit) (bool, error) {
		if len(commit.Message) > 10 {
			return false, nil
		}

		return true, nil
	}
}

// All commits on the main branch for github flow should be merged in,
// meaning that they have two parents(the main branch and the feature branch).
func TwoParentsCommitDetect() CommitDetect {
	return func(commit *model.Commit) (bool, error) {
		if len(commit.ParentHashes) >= 2 {
			return true, nil
		}

		return false, nil
	}
}
