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
			// TODO: Force push types
			// forcePush := detector.NewCommitCacheDetector(detector.ForcePushDetect())
			unresolved := detector.NewCommitDetector(detector.UnresolvedDetect())

			g := options.GetGradingAlgorithm(settings.GradingAlgorithm, settings.ThresholdSettings)

			return "General", analysis.DetectorMarker(m,
				[]detector.Detector{unresolved},
				m.Contribution.CommitCountMap,
				1, g)
		},
	)

	return analyzer
}
