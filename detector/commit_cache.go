package detector

import (
	"encoding/hex"
	"fmt"

	"github.com/Git-Gopher/go-gopher/cache"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/Git-Gopher/go-gopher/violation"
)

type CommitCacheDetect func(c *common, email string, current *cache.Cache, cache *cache.Cache) (bool, []violation.Violation, error)

type CommitCacheDetector struct {
	violated   int
	found      int
	total      int
	violations []violation.Violation

	detect CommitCacheDetect
}

// TODO: We should change this to the enriched model.
func (cd *CommitCacheDetector) Run(owner, repo, email string, current *cache.Cache, cache []*cache.Cache) error {
	// Struct should be reset before each run, incase we are running it with a different model.
	cd.violated = 0
	cd.found = 0
	cd.total = 0
	cd.violations = make([]violation.Violation, 0)
	c := common{owner: owner, repo: repo}

	for _, ca := range cache {
		found, vlns, err := cd.detect(&c, email, current, ca)
		if err != nil {
			return fmt.Errorf("Error running cache detector: %w", err)
		}

		if found {
			cd.found++
		}
		cd.violations = append(cd.violations, vlns...)
		cd.total++
	}

	return nil
}

func (cd *CommitCacheDetector) Result() (int, int, int, []violation.Violation) {
	return cd.violated, cd.found, cd.total, cd.violations
}

func NewCommitCacheDetector(detect CommitCacheDetect) *CommitCacheDetector {
	return &CommitCacheDetector{
		violated:   0,
		found:      0,
		total:      0,
		violations: make([]violation.Violation, 0),
		detect:     detect,
	}
}

// GithubWorklow: Force pushes are not allowed.
func ForcePushDetect() CommitCacheDetect {
	return func(c *common, email string, current *cache.Cache, cache *cache.Cache) (bool, []violation.Violation, error) {
		lhs := make([]utils.Commit, 0)
		for _, cuh := range current.Hashes {
			for _, cah := range cache.Hashes {
				if cuh == cah {
					return false, nil, nil
				}
			}
			// Hash not found in cache
			lh := utils.Commit{
				Hash: hex.EncodeToString(cuh.ToByte()),
				GitHubLink: utils.GitHubLink{
					Owner: c.owner,
					Repo:  c.repo,
				},
			}
			lhs = append(lhs, lh)
		}

		violations := [1]violation.Violation{violation.NewForcePushViolation(lhs, email)}

		return true, violations[:], nil
	}
}
