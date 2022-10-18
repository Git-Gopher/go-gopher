package markers

import (
	"github.com/Git-Gopher/go-gopher/assess/markers/analysis"
	"github.com/Git-Gopher/go-gopher/assess/options"
	"github.com/Git-Gopher/go-gopher/detector"
)

const branchingName = "branch"

func NewBranching(settings *options.BranchingSettings) *analysis.Analyzer {
	analyzer := analysis.NewAnalyzer(
		branchingName,
		"Branching marker",
		func(m analysis.MarkerCtx) (string, []analysis.Mark) {
			stale := detector.NewBranchDetector(detector.StaleBranchDetect())
			consistent := detector.NewBranchCompareDetector(detector.BranchNameConsistencyDetect())
			feature := detector.NewFeatureBranchDetector("FeatureBranchDetector")

			g := options.GetGradingAlgorithm(settings.GradingAlgorithm, settings.ThresholdSettings)

			return "Branching", analysis.DetectorMarker(m,
				[]detector.Detector{consistent, stale, feature},
				m.Contribution.CommitCountMap,
				3, g)
		},
	)

	return analyzer
}
