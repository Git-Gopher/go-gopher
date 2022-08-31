package detector

import (
	"fmt"
	"time"

	"github.com/Git-Gopher/go-gopher/markup"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/Git-Gopher/go-gopher/violation"
	log "github.com/sirupsen/logrus"
)

type BranchDetect func(c *common, branch *local.Branch) (bool, violation.Violation, error)

// BranchDetector is used to run a detector on each branch metadata.
type BranchDetector struct {
	name       string
	violated   int
	found      int
	total      int
	violations []violation.Violation

	detect BranchDetect
}

func NewBranchDetector(name string, detect BranchDetect) *BranchDetector {
	return &BranchDetector{
		name:       name,
		violated:   0,
		found:      0,
		total:      0,
		violations: make([]violation.Violation, 0),
		detect:     detect,
	}
}

func (bd *BranchDetector) Run(em *enriched.EnrichedModel) error {
	if em == nil {
		return nil
	}

	bd.violated = 0
	bd.found = 0
	bd.total = 0
	bd.violations = make([]violation.Violation, 0)
	c, err := NewCommon(em)
	if err != nil {
		log.Fatalf("could not create common: %v", err)
	}

	for _, b := range em.Branches {
		b := b
		detected, violation, err := bd.detect(c, &b)
		if err != nil {
			return fmt.Errorf("failed to run branch detector: %w", err)
		}
		if err != nil {
			return err
		}
		if detected {
			bd.found++
		}
		if violation != nil {
			bd.violations = append(bd.violations, violation)
		}
		bd.total++
	}

	return nil
}

func (b *BranchDetector) Result() (int, int, int, []violation.Violation) {
	return b.violated, b.found, b.total, b.violations
}

func (b *BranchDetector) Name() string {
	return b.name
}

// GithubWorklow: Branches are considered stale after three months.
func StaleBranchDetect() (string, BranchDetect) {
	// NOTE: Course is currently using 2 weeks as stale branch time because of the short life of the project.
	// We should make this configurable in the future!
	secondsInWeek := 604800
	staleBranchTime := time.Hour * 24 * 14 // 14 days

	// secondsInMonth := 2600640

	return "StaleBranchDetect", func(c *common, branch *local.Branch) (bool, violation.Violation, error) {
		if time.Since(branch.Head.Committer.When) > staleBranchTime {
			email := branch.Head.Committer.Email
			monthsSince := utils.RoundTime(time.Since(branch.Head.Committer.When).Seconds() / float64(secondsInWeek))

			return true, violation.NewStaleBranchViolation(
				markup.Branch{
					Name: branch.Name,
					GitHubLink: markup.GitHubLink{
						Owner: c.owner,
						Repo:  c.repo,
					},
				},
				time.Duration(monthsSince),
				email,
				c.IsCurrentBranch(branch.Name),
			), nil
		}

		return false, nil, nil
	}
}
