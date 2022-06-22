package detector

import (
	"github.com/Git-Gopher/go-gopher/model/github"
	"github.com/Git-Gopher/go-gopher/model/local"
)

type PullRequestDetect func(pullRequest *github.PullRequest) (bool, error)

// XXX: violated, found, total should be contained within a struct and then added to this instead as a composite struct.
type PullRequestDetector struct {
	violated int
	found    int
	total    int

	detect PullRequestDetect
}

func (pd *PullRequestDetector) Result() (int, int, int) {
	return pd.violated, pd.found, pd.total
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
func (pd *PullRequestDetector) Run(model *local.GitModel) error {
	return nil
}

func (pd *PullRequestDetector) Run2(model *github.GithubModel) error {
	for _, pr := range model.PullRequests {
		pr := pr
		detected, err := pd.detect(pr)
		pd.total++
		if err != nil {
			return err
		}
		if detected {
			pd.found++
		}
	}

	return nil
}

// Github Workflow: Pull requests must have at least one associated issue.
func PullRequestIssueDetector() PullRequestDetect {
	return func(pullRequest *github.PullRequest) (bool, error) {
		if len(pullRequest.ClosingIssues) == 0 {
			return false, nil
		}

		return true, nil
	}
}
