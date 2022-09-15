package detector

import (
	"github.com/Git-Gopher/go-gopher/markup"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/remote"
	"github.com/Git-Gopher/go-gopher/violation"
	log "github.com/sirupsen/logrus"
)

type PullRequestDetect func(c *common, pullRequest *remote.PullRequest) (bool, violation.Violation, error)

// XXX: violated, found, total should be contained within a struct and then added to this instead as a composite struct.
type PullRequestDetector struct {
	name       string
	violated   int
	found      int
	total      int
	violations []violation.Violation

	detect PullRequestDetect
}

func NewPullRequestDetector(name string, detect PullRequestDetect) *PullRequestDetector {
	return &PullRequestDetector{
		name:     name,
		violated: 0,
		found:    0,
		total:    0,
		detect:   detect,
	}
}

func (pd *PullRequestDetector) Run(em *enriched.EnrichedModel) error {
	c, err := NewCommon(em)
	if err != nil {
		log.Printf("could not create common: %v", err)
	}
	for _, pr := range em.PullRequests {
		pr := pr
		detected, violation, err := pd.detect(c, pr)
		pd.total++
		if err != nil {
			return err
		}
		if detected {
			pd.found++
		}
		if violation != nil {
			pd.violations = append(pd.violations, violation)
		}
	}

	return nil
}

func (prd *PullRequestDetector) Result() (int, int, int, []violation.Violation) {
	return prd.violated, prd.found, prd.total, prd.violations
}

func (prd *PullRequestDetector) Name() string {
	return prd.name
}

// Github Workflow: Pull requests must have at least one associated issue.
func PullRequestIssueDetector() (string, PullRequestDetect) {
	return "PullRequestIssueDetector", func(c *common, pr *remote.PullRequest) (bool, violation.Violation, error) {
		if len(pr.ClosingIssues) == 0 {
			if pr.CreatedAt == nil {
				return false, nil, violation.ErrCreatedTimePullRequest
			}

			return true, violation.NewLinkedIssueViolation(
				markup.PR{
					Number: pr.Number,
					GitHubLink: markup.GitHubLink{
						Owner: c.owner,
						Repo:  c.repo,
					},
				},
				c.IsCurrentPR(pr),
				*pr.CreatedAt,
			), nil
		}

		return false, nil, nil
	}
}

func PullRequestApprovalDetector() (string, PullRequestDetect) {
	return "PullRequestApprovalDetector", func(c *common, pr *remote.PullRequest) (bool, violation.Violation, error) {
		// Ignore unmerged pull requests.
		if !pr.Merged {
			return false, nil, nil
		}

		if pr.ReviewDecision != "APPROVED" {
			// Pull request must have closed time if merged.
			if pr.ClosedAt == nil {
				return false, nil, violation.ErrClosedTimePullRequest
			}

			return true, violation.NewApprovalViolation(
				markup.PR{
					Number: pr.Number,
					GitHubLink: markup.GitHubLink{
						Owner: c.owner,
						Repo:  c.repo,
					},
				},
				c.IsCurrentPR(pr),
				*pr.ClosedAt,
			), nil
		}

		return false, nil, nil
	}
}

// All reviews threads should be marked as resolved before merging.
func PullRequestReviewThreadDetector() (string, PullRequestDetect) {
	return "PullRequestReviewThreadDetector", func(c *common, pr *remote.PullRequest) (bool, violation.Violation, error) {
		// Ignore open and unmerged pull requests.
		if !pr.Closed || !pr.Merged {
			return false, nil, nil
		}

		for _, thread := range pr.ReviewThreads {
			if !thread.IsResolved {
				if pr.ClosedAt == nil {
					return false, nil, violation.ErrClosedTimePullRequest
				}

				return true, violation.NewUnresolvedConversationViolation(
					markup.PR{
						Number: pr.Number,
						GitHubLink: markup.GitHubLink{
							Owner: c.owner,
							Repo:  c.repo,
						},
					},
					c.IsCurrentPR(pr),
					*pr.ClosedAt,
				), nil
			}
		}

		return false, nil, nil
	}
}
