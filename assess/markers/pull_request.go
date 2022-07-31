package markers

import (
	"github.com/Git-Gopher/go-gopher/assess/markers/analysis"
	"github.com/Git-Gopher/go-gopher/assess/options"
	"github.com/Git-Gopher/go-gopher/detector"
)

const pullRequestName = "pull-request"

func NewPullRequest(settings *options.PullRequestSettings) *analysis.Analyzer {
	analyzer := analysis.NewAnalyzer(
		pullRequestName,
		"PullRequest marker",
		func(m analysis.MarkerCtx) (string, []analysis.Mark) {
			approved := detector.NewPullRequestDetector(detector.PullRequestApprovalDetector())
			resolved := detector.NewPullRequestDetector(detector.PullRequestReviewThreadDetector())
			issue := detector.NewPullRequestDetector(detector.PullRequestIssueDetector())

			g := options.GetGradingAlgorithm(settings.GradingAlgorithm, settings.ThresholdSettings)
			return "PullRequest", analysis.DetectorMarker(m,
				[]detector.Detector{approved, resolved, issue},
				m.Contribution.CommitCountMap,
				3, g)
		},
	)

	return analyzer
}
