package detector

import (
	"fmt"

	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/violation"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	log "github.com/sirupsen/logrus"
)

// XXX: This detector should be avoided unless creating detectors that work on
// multiple domains (commits, tags, etc) which are unlikely to be reused in the future.
// This detector has the ability to work on a lot of data due to working on the root
// repository which could lead to high complexity and runtime.
type EnrichedModelDetect func(c *common, em *enriched.EnrichedModel) (int, []violation.Violation, error)

type EnrichedModelDetector struct {
	name       string
	violated   int
	found      int
	total      int
	violations []violation.Violation

	detect EnrichedModelDetect
}

func NewEnrichedModelDetector(name string, detect EnrichedModelDetect) *EnrichedModelDetector {
	return &EnrichedModelDetector{
		name:       name,
		violated:   0,
		found:      0,
		total:      0,
		violations: make([]violation.Violation, 0),
		detect:     detect,
	}
}

func (b *EnrichedModelDetector) Run(em *enriched.EnrichedModel) error {
	if em == nil {
		return nil
	}

	b.violated = 0
	b.found = 0
	b.total = 0
	b.violations = make([]violation.Violation, 0)
	c, err := NewCommon(em)
	if err != nil {
		log.Printf("could not create common: %v", err)
	}

	b.found, b.violations, err = b.detect(c, em)
	if err != nil {
		return fmt.Errorf("error EnrichedModelDetector: %w", err)
	}

	return nil
}

func (emd *EnrichedModelDetector) Result() (int, int, int, []violation.Violation) {
	return emd.violated, emd.found, emd.total, emd.violations
}

func (emd *EnrichedModelDetector) Name() string {
	return emd.name
}

func TagMainCommitsDetect() (string, EnrichedModelDetect) {
	return "TagMainCommitsDetect", func(c *common, em *enriched.EnrichedModel) (int, []violation.Violation, error) {
		// XXX: replace with main detection algorithm.
		var primaryBranchCommitHashes []local.Hash
		primaryBranch := "main"
		branchHeadRefName := fmt.Sprintf("refs/remotes/origin/%s", primaryBranch)

		// Fetch commits that exist on primary branch.
		primaryRef, err := em.Repository.Reference(plumbing.ReferenceName(branchHeadRefName), false)
		if err != nil {
			return 0, nil, fmt.Errorf("could not fetch main branch reference: %w", err)
		}

		primaryIter, err := em.Repository.Log(&git.LogOptions{
			From:  primaryRef.Hash(),
			Order: git.LogOrderCommitterTime,
		})
		if err != nil {
			return 0, nil, fmt.Errorf("error creating commit iter for branch: %w", err)
		}

		if err = primaryIter.ForEach(func(c *object.Commit) error {
			primaryBranchCommitHashes = append(primaryBranchCommitHashes, local.Hash(c.Hash))

			return nil
		}); err != nil {
			return 0, nil, fmt.Errorf("error folding primary branch commits: %w", err)
		}

		// Move tags to set
		tags := make(map[local.Hash]struct{})
		for _, t := range em.Tags {
			tags[t.Hash] = struct{}{}
		}

		// Proportion of tags that exist for main commits.
		bad := 0
		good := 0
		for _, c := range primaryBranchCommitHashes {
			if _, ok := tags[c]; !ok {
				bad += 1
			} else {
				good += 1
			}
		}

		log.Printf("number commits that have tags: %v\n", good)
		log.Printf("number commits that don't have tags : %v\n", bad)

		// This interface that we are using might be slightly limiting...
		// I can't really access the total directly from this detector.
		// Plus this would mess up a total for other runs.
		return 0, nil, nil
	}
}
