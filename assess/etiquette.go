package assess

import (
	"github.com/Git-Gopher/go-gopher/detector"
)

// TODO: Make the configuration for each of these markers configurable

func Commit(m MarkerCtx) (string, []Mark) {
	atomicity := detector.NewCommitDistanceDetector(detector.DiffDistanceCalculation())
	binaries := detector.NewCommitDetector(detector.BinaryDetect())
	empty := detector.NewCommitDetector(detector.EmptyCommitDetect())

	return "Commit", DetectorMarker(m, []detector.Detector{atomicity, binaries, empty}, m.Contribution.CommitCountMap)
}

func CommitMessage(m MarkerCtx) (string, []Mark) {
	diff := detector.NewCommitDetector(detector.DiffMatchesMessageDetect())
	short := detector.NewCommitDetector(detector.ShortCommitMessageDetect())

	return "CommitMessage", DetectorMarker(m, []detector.Detector{diff, short}, m.Contribution.CommitCountMap)
}

func Branching(m MarkerCtx) (string, []Mark) {
	stale := detector.NewBranchDetector(detector.StaleBranchDetect())
	consistent := detector.NewBranchCompareDetector(detector.BranchNameConsistencyDetect())
	feature := detector.NewFeatureBranchDetector("FeatureBranchDetector")

	return "Branching", DetectorMarker(m, []detector.Detector{consistent, stale, feature}, m.Contribution.CommitCountMap)
}

func PullRequest(m MarkerCtx) (string, []Mark) {
	approved := detector.NewPullRequestDetector(detector.PullRequestApprovalDetector())
	resolved := detector.NewPullRequestDetector(detector.PullRequestReviewThreadDetector())
	issue := detector.NewPullRequestDetector(detector.PullRequestIssueDetector())

	return "PullRequest", DetectorMarker(m, []detector.Detector{approved, resolved, issue}, m.Contribution.CommitCountMap)
}

func General(m MarkerCtx) (string, []Mark) {
	// TODO: Force push types
	// forcePush := detector.NewCommitCacheDetector(detector.ForcePushDetect())
	unresolved := detector.NewCommitDetector(detector.UnresolvedDetect())

	return "General", DetectorMarker(m, []detector.Detector{unresolved}, m.Contribution.CommitCountMap)
}
