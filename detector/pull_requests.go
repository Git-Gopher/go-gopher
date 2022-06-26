package detector

import (
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/github"
	"github.com/Git-Gopher/go-gopher/violation"
)

type PullRequestDetect func(pullRequest *github.PullRequest) (bool, violation.Violation, error)

// XXX: violated, found, total should be contained within a struct and then added to this instead as a composite struct.
type PullRequestDetector struct {
	violated   int
	found      int
	total      int
	violations []violation.Violation

	detect PullRequestDetect
}

func (prd *PullRequestDetector) Result() (int, int, int, []violation.Violation) {
	return prd.violated, prd.found, prd.total, prd.violations
}

func NewPullRequestDetector(detect PullRequestDetect) *PullRequestDetector {
	return &PullRequestDetector{
		violated: 0,
		found:    0,
		total:    0,
		detect:   detect,
	}
}

// TODO: We should change this to the enriched model.
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

// Github Workflow: Pull requests must have at least one associated issue.
func PullRequestIssueDetector() PullRequestDetect {
	return func(pullRequest *github.PullRequest) (bool, violation.Violation, error) {
		if len(pullRequest.ClosingIssues) == 0 {
			return false, nil, nil
		}

		return true, nil, nil
	}
}

func PullRequestApprovalDetector() PullRequestDetect {
	return func(pullRequest *github.PullRequest) (bool, violation.Violation, error) {
		// Ignore unmerged pull requests.
		if pullRequest.Merged {
			return true, nil, nil
		}

		// XXX: Create enum
		if pullRequest.ReviewDecision != "APPROVED" {
			return false, nil, nil
		}

		return true, nil, nil
	}
}

// All reviews threads should be marked as resolved before merging.
func PullRequestReviewThreadDetector() PullRequestDetect {
	return func(pullRequest *github.PullRequest) (bool, violation.Violation, error) {
		// Ignore unmerged pull requests.
		if pullRequest.Merged {
			return true, nil, nil
		}

		for _, thread := range pullRequest.ReviewThreads {
			if !thread.IsResolved {
				return false, nil, nil
			}
		}

		return true, nil, nil
	}
}
