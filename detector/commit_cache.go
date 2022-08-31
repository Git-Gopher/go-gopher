package detector

import (
	"fmt"

	"github.com/Git-Gopher/go-gopher/cache"
	"github.com/Git-Gopher/go-gopher/markup"
	"github.com/Git-Gopher/go-gopher/violation"
)

type CommitCacheDetect func(
	c *common,
	email string,
	current *cache.Cache,
	previous []*cache.Cache,
) (bool, []violation.Violation, error)

type CommitCacheDetector struct {
	name       string
	violated   int
	found      int
	total      int
	violations []violation.Violation

	detect CommitCacheDetect
}

func NewCommitCacheDetector(name string, detect CommitCacheDetect) *CommitCacheDetector {
	return &CommitCacheDetector{
		name:       name,
		violated:   0,
		found:      0,
		total:      0,
		violations: make([]violation.Violation, 0),
		detect:     detect,
	}
}

func (cd *CommitCacheDetector) Run(owner, repo, email string, current *cache.Cache, previous []*cache.Cache) error {
	// Struct should be reset before each run, incase we are running it with a different model.
	cd.violated = 0
	cd.found = 0
	cd.total = 0
	cd.violations = make([]violation.Violation, 0)
	c := common{owner: owner, repo: repo}

	found, vs, err := cd.detect(&c, email, current, previous)
	if err != nil {
		return fmt.Errorf("Error running cache detector: %w", err)
	}

	if found {
		cd.found++
	}

	cd.violations = append(cd.violations, vs...)
	cd.total++

	return nil
}

func (cd *CommitCacheDetector) Result() (int, int, int, []violation.Violation) {
	return cd.violated, cd.found, cd.total, cd.violations
}

func (cd *CommitCacheDetector) Name() string {
	return cd.name
}

// GithubWorklow: Force pushes are not allowed.
func ForcePushDetect() (string, CommitCacheDetect) {
	return "ForcePushDetect",
		func(c *common, email string, current *cache.Cache, previous []*cache.Cache) (bool, []violation.Violation, error) {
			missing := make([]markup.Commit, 0)
			for _, pc := range previous {
				hashes := make([]string, 0, len(pc.Hashes))
				for k := range pc.Hashes {
					hashes = append(hashes, k)
				}

				for _, h := range hashes {
					if _, ok := current.Hashes[h]; !ok {
						missing = append(missing,
							markup.Commit{
								Hash: h,
								GitHubLink: markup.GitHubLink{
									Owner: c.owner,
									Repo:  c.repo,
								},
							},
						)
					}
				}
			}

			if len(missing) == 0 {
				return false, nil, nil
			}

			// XXX: Force pushes will always show
			violations := [1]violation.Violation{
				violation.NewForcePushViolation(missing, email, current.Created, true),
			}

			return true, violations[:], nil
		}
}
