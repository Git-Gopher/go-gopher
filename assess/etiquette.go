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

// E-A1: Commit message:
// Commit messages should be descriptive and concise.
// Commit messages should not be too short.
func CommitMessage(m MarkerCtx) (string, []Mark) {
	diff := detector.NewCommitDetector(detector.DiffMatchesMessageDetect())

	short := detector.NewCommitDetector(detector.ShortCommitMessageDetect())

	diffMarker := DetectorMarker(m, diff, m.Contribution.CommitCountMap)

	shortMarker := DetectorMarker(m, short, m.Contribution.CommitCountMap)

	diffMarkMap := make(map[string]Mark)
	for _, mark := range diffMarker {
		diffMarkMap[mark.Username] = mark
	}

	for _, mark := range shortMarker {
		if _, ok := diffMarkMap[mark.Username]; !ok {
			diffMarkMap[mark.Username] = mark

			continue
		}

		diffMark := diffMarkMap[mark.Username]
		diffMark.Violations = append(diffMark.Violations, mark.Violations...)
		diffMark.Total = diffMark.Total + mark.Total
		diffMarkMap[mark.Username] = diffMark
	}

	allMarks := make([]Mark, 0, len(diffMarkMap))
	for _, mark := range diffMarkMap {
		allMarks = append(allMarks, mark)
	}

	return "CommitMessage", allMarks
}

// E-A1: Generated files (low priority); we can detect binaries etc.
func GeneratedFiles(m MarkerCtx) (string, []Mark) {
	// Not Implemented.
	return "GeneratedFiles", []Mark{}
}

// E-A2: Regular branch names: Use regular branch name and branch name prefixes that
// accurately represent the work that the branch contains.
func RegularBranchNames(m MarkerCtx) (string, []Mark) {
	d := detector.NewBranchCompareDetector(detector.BranchNameConsistencyDetect())

	return "RegularBranchNames", DetectorMarker(m, d, m.Contribution.BranchCountMap)
}
