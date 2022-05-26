package detector

import "github.com/Git-Gopher/go-gopher/model"

type CommitDetect func(commit *model.Commit) (bool, error)

type CommitDetector struct {
	violated int
	found    int
	total    int

	detect CommitDetect
}

// XXX: Run needs to be moved up to the detector inteface and made generic on the object it recieves, or we can just set it to be the entire git tree so that we don't have to think about it
func (c *CommitDetector) Run(commit *model.Commit) error {
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
	return func(commit *model.Commit) (bool, error) {
		if len(commit.Message) > 10 {
			return false, nil
		}

		return true, nil
	}
}
