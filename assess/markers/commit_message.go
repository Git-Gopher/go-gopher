package markers

import (
	"github.com/Git-Gopher/go-gopher/assess/markers/analysis"
	"github.com/Git-Gopher/go-gopher/assess/options"
	"github.com/Git-Gopher/go-gopher/detector"
)

const commitMessageName = "commit-message"

func NewCommitMessage(settings *options.CommitMessageSetting) *analysis.Analyzer {
	analyzer := analysis.NewAnalyzer(
		commitMessageName,
		"Commit Message marker",
		func(m analysis.MarkerCtx) (string, []analysis.Mark) {
			diff := detector.NewCommitDetector(detector.DiffMatchesMessageDetect())
			short := detector.NewCommitDetector(detector.ShortCommitMessageDetect())

			g := options.GetGradingAlgorithm(settings.GradingAlgorithm, settings.ThresholdSettings)

			return "CommitMessage", analysis.DetectorMarker(m,
				[]detector.Detector{diff, short},
				m.Contribution.CommitCountMap,
				2, g)
		},
	)

	return analyzer
}
