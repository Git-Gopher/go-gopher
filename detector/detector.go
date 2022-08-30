package detector

import (
	"fmt"
	"log"

	"github.com/Git-Gopher/go-gopher/cache"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/model/remote"
	"github.com/Git-Gopher/go-gopher/violation"
)

var (
	ErrNotImplemented = fmt.Errorf("Not implemented")

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

// Create a common object from the enriched model.
func NewCommon(em *enriched.EnrichedModel) (*common, error) {
	if commonMemo == nil {
		var mergingCommits []local.Hash
		currentPR, err := em.FindCurrentPR()
		if err != nil {
			log.Printf("could not find current PR: %v", err)
		} else {
			mergingCommits, err = em.FindMergingCommits(currentPR)
			if err != nil {
				log.Printf("could not find merging commits: %v", err)
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
