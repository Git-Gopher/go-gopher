package detector

import (
	"fmt"

	"github.com/Git-Gopher/go-gopher/cache"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/model/remote"
	"github.com/Git-Gopher/go-gopher/violation"
	log "github.com/sirupsen/logrus"
)

var (
	ErrNotImplemented = fmt.Errorf("not implemented")

	commonMemo *common
)

// common - common variables that are shared with all detectors.
type common struct {
	// Owner of repository.
	owner string
	// Name of the repository.
	repo string
	// Commits that are going to be merged into target branch from source branch.
	mergingCommits []local.Hash
	// Current pull request.
	PR *remote.PullRequest
}

// Checks if a commit relates to the current feedback comment.
func (c *common) IsCurrentCommit(h local.Hash) bool {
	// In the case where there is no current branch/pr, default to report all.
	if c.mergingCommits == nil {
		return true
	}

	for _, v := range c.mergingCommits {
		if v == h {
			return true
		}
	}

	return false
}

func (c *common) IsCurrentPR(pr *remote.PullRequest) bool {
	// In the case where there is no current branch/pr, default to report all.
	if c.PR == nil {
		log.Warn("no pr found")

		return true
	}

	if pr.Number == c.PR.Number {
		log.Warnf("current pr number %d", c.PR.Number)

		return true
	}

	return false
}

func (c *common) IsCurrentBranch(branchName string) bool {
	// In the case where there is no current branch/pr, default to report all.
	if c.PR == nil {
		return true
	}

	if c.PR.HeadRefName == branchName {
		return true
	}

	return false
}

// Create a common object from the enriched model.
func NewCommon(em *enriched.EnrichedModel) (*common, error) {
	if commonMemo == nil {
		var mergingCommits []local.Hash
		currentPR, err := em.FindCurrentPR()
		if err != nil {
			log.Warn("unable to find current PR")
		} else {
			mergingCommits, err = em.FindMergingCommits(currentPR)
			if err != nil {
				log.Warn("unable to find merging commits")
			}
		}

		commonMemo = &common{
			owner:          em.Owner,
			repo:           em.Name,
			PR:             currentPR,
			mergingCommits: mergingCommits,
		}
	}

	return commonMemo, nil
}

type Detector interface {
	Run(model *enriched.EnrichedModel) error
	Result() (violated, count, total int, violations []violation.Violation)
	Name() string
}

type CacheDetector interface {
	Run(owner string, repo string, email string, current *cache.Cache, previous []*cache.Cache) error
	Result() (violated, count, total int, violations []violation.Violation)
	Name() string
}
