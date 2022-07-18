package assess

import (
	"github.com/Git-Gopher/go-gopher/detector"
)

// W-A2: How it should be: Features are made on separate branches and then merged
// into main via a pull request instead of directly committing them to the main branch.
func FeatureBranching(m MarkerCtx) (string, []Mark) {
	d := detector.NewFeatureBranchDetector()

	return "FeatureBranching", DetectorMarker(m, d, m.Contribution.MergeCountMap)
}

// W-A2: Pull requests are created for features and atleast one person in your team
// should review your code before merging it into main.
func PullRequestReview(m MarkerCtx) (string, []Mark) {
	d := detector.NewPullRequestDetector(detector.PullRequestReviewThreadDetector())

	return "PullRequestReview", DetectorMarker(m, d, m.Contribution.PullRequestCountMap)
}
