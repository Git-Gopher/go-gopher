package detector

import (
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/remote"
	"github.com/Git-Gopher/go-gopher/violation"
)

type PullRequestDetect func(pullRequest *remote.PullRequest) (bool, violation.Violation, error)

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

func (pd *PullRequestDetector) Run(model *enriched.EnrichedModel) error {
	for _, pr := range model.PullRequests {
		pr := pr
		detected, violation, err := pd.detect(pr)
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
	return "PullRequestIssueDetector", func(pr *remote.PullRequest) (bool, violation.Violation, error) {
		if len(pr.ClosingIssues) == 0 {
			return true, violation.NewLinkedIssueViolation(pr.Url), nil
		}

		return false, nil, nil
	}
}

func PullRequestApprovalDetector() (string, PullRequestDetect) {
	return "PullRequestApprovalDetector", func(pr *remote.PullRequest) (bool, violation.Violation, error) {
		// Ignore unmerged pull requests.
		if !pr.Merged {
			return false, nil, nil
		}

		if pr.ReviewDecision != "APPROVED" {
			return true, violation.NewApprovalViolation(pr.Url), nil
		}

		return false, nil, nil
	}
}

// All reviews threads should be marked as resolved before merging.
func PullRequestReviewThreadDetector() (string, PullRequestDetect) {
	return "PullRequestReviewThreadDetector", func(pr *remote.PullRequest) (bool, violation.Violation, error) {
		// Ignore unmerged pull requests.
		if pr.Merged {
			return false, nil, nil
		}

		for _, thread := range pr.ReviewThreads {
			if !thread.IsResolved {
				return true, violation.NewUnresolvedConversationViolation(pr.Url), nil
			}
		}

		return false, nil, nil
	}
}
