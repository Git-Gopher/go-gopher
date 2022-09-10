package detector

import (
	"fmt"

	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/violation"
)

type CommitMainTaggedDetector struct {
	name     string
	violated int
	found    int
	total    int
}

func NewCommitMainTaggedDetector(name string) *CommitMainTaggedDetector {
	return &CommitMainTaggedDetector{
		name:     name,
		violated: 0,
		found:    0,
		total:    0,
	}
}

func (cd *CommitMainTaggedDetector) Run(em *enriched.EnrichedModel) error {
	if em == nil {
		return nil
	}

	// Struct should be reset before each run, incase we are running it with a different model.
	cd.violated = 0
	cd.found = 0
	cd.total = 0

	// TODO: Replace with primary branch detector
	commits, err := em.CommitsOnBranch("main")
	if err != nil {
		return fmt.Errorf("could not fetch commits on branch \"main\": %w", err)
	}

	// Memo tags
	tags := make(map[local.Hash]struct{})
	for _, t := range em.Tags {
		if _, ok := tags[t.Hash]; !ok {
			tags[t.Hash] = struct{}{}
		}
	}

	for _, c := range commits {
		c := c
		// Commit doesn't have associated tag.
		if _, ok := tags[c]; !ok {
			cd.violated++
		} else {
			cd.found++
		}
	}

	cd.total = len(commits)

	return nil
}

func (cd *CommitMainTaggedDetector) Result() (int, int, int, []violation.Violation) {
	return cd.violated, cd.found, cd.total, nil
}

func (cd *CommitMainTaggedDetector) Name() string {
	return cd.name
}
