package markers

import (
	"github.com/Git-Gopher/go-gopher/assess/markers/analysis"
	"github.com/Git-Gopher/go-gopher/assess/options"
	"github.com/Git-Gopher/go-gopher/detector"
)

const commitName = "commit"

func NewCommit(settings *options.CommitSettings) *analysis.Analyzer {
	analyzer := analysis.NewAnalyzer(
		commitName,
		"Commit marker",
		func(m analysis.MarkerCtx) (string, []analysis.Mark) {
			atomicity := detector.NewCommitDistanceDetector(detector.DiffDistanceCalculation())
			binaries := detector.NewCommitDetector(detector.BinaryDetect())
			empty := detector.NewCommitDetector(detector.EmptyCommitDetect())

			g := options.GetGradingAlgorithm(settings.GradingAlgorithm, settings.ThresholdSettings)
			return "Commit", analysis.DetectorMarker(
				m,
				[]detector.Detector{atomicity, binaries, empty},
				m.Contribution.CommitCountMap,
				3, g)
		},
	)

	return analyzer
}
