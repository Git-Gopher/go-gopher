package assess

import (
	"github.com/Git-Gopher/go-gopher/assess/assessors"
	"github.com/Git-Gopher/go-gopher/assess/markers/analysis"
	"github.com/Git-Gopher/go-gopher/assess/options"
	log "github.com/sirupsen/logrus"
)

func LoadAnalyzer(option *options.Options) []*analysis.Analyzer {
	logger := log.New()

	m := assessors.NewManager(logger, option)

	es := assessors.NewEnabledSet(m, option)

	analyzerMap := es.GetEnabledMarkersMap()
	analyzers := make([]*analysis.Analyzer, 0, len(analyzerMap))
	for _, analyzer := range analyzerMap {
		analyzers = append(analyzers, analyzer)
	}

	return analyzers
}
