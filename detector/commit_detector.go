package detector

import "github.com/go-git/go-git/v5/plumbing/object"

type Commit struct {
	Commit *object.Commit // remove future

	Parent  string
	Hash    string
	Message string
	Content string
}

type CommitDetect func(commit *Commit) (bool, error)

type CommitDetector struct {
	violated int
	found    int
	total    int

	detect CommitDetect
}

func (c *CommitDetector) Run(commit *Commit) error {
	c.total++
	detected, err := c.detect(commit)
	if err != nil {
		return err
	}
	if detected {
		c.found++
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
	return func(commit *Commit) (bool, error) {
		if len(commit.Message) > 10 {
			return false, nil
		}

		return true, nil
	}
}
