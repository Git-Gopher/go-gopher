package assess

import (
	"github.com/Git-Gopher/go-gopher/detector"
)

// E-A1: Diff Distance/Atomicity:
// Commit diffs should be atomic and only encapsulate a single change to the system,
// commits that include multiple different conceuptal changes should instead be
// multiple different commits.
func Atomicity(m MarkerCtx) (string, []Mark) {
	d := detector.NewCommitDistanceDetector(detector.DiffDistanceCalculation())

	return "Atomicity", DetectorMarker(m, d, m.Contribution.CommitCountMap)
}

// E-A1: Commit messages are descriptive and relate to the contents of the commit
// diff that contains the change made to the code (eg: "fix: index offset incorrect
// starting value bug".
func DescriptiveCommit(m MarkerCtx) (string, []Mark) {
	d := detector.NewCommitDetector(detector.DiffMatchesMessageDetect())

	return "DescriptiveCommit", DetectorMarker(m, d, m.Contribution.CommitCountMap)
}

// E-A1: Generated files (low priority); we can detect binaries etc.
func GeneratedFiles(m MarkerCtx) (string, []Mark) {
	// Not Implemented.
	return "GeneratedFiles", []Mark{}
}

// E-A2: Regular branch names: Use regular branch name and branch name prefixes that
// accurately represent the work that the branch contains.
func RegularBranchNames(m MarkerCtx) (string, []Mark) {
	d := detector.NewBranchCompareDetector(detector.NewBranchNameConsistencyDetect())

	return "RegularBranchNames", DetectorMarker(m, d, m.Contribution.BranchCountMap)
}
