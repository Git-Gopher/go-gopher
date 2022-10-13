package markers

import (
	"github.com/Git-Gopher/go-gopher/assess/markers/analysis"
	"github.com/Git-Gopher/go-gopher/assess/options"
	"github.com/Git-Gopher/go-gopher/detector"
)

const generalName = "general"

func NewGeneral(settings *options.GeneralSettings) *analysis.Analyzer {
	analyzer := analysis.NewAnalyzer(
		generalName,
		"General marker",
		func(m analysis.MarkerCtx) (string, []analysis.Mark) {
			unresolved := detector.NewCommitDetector(detector.UnresolvedDetect())
			featureBranch := detector.NewFeatureBranchDetector("FeatureBranchDetector")

			g := options.GetGradingAlgorithm(settings.GradingAlgorithm, settings.ThresholdSettings)

			return "General", analysis.DetectorMarker(m,
				[]detector.Detector{unresolved, featureBranch},
				m.Contribution.CommitCountMap,
				1, g)
		},
	)

	return analyzer
}
