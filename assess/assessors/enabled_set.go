package assessors

import (
	"github.com/Git-Gopher/go-gopher/assess/markers/analysis"
	"github.com/Git-Gopher/go-gopher/assess/options"
)

type EnabledSet struct {
	m   *Manager
	opt *options.Options
}

func NewEnabledSet(m *Manager, opt *options.Options) *EnabledSet {
	return &EnabledSet{
		m:   m,
		opt: opt,
	}
}

func (es EnabledSet) build(marker *options.Markers) map[string]*analysis.Analyzer {
	resultMarkersSet := make(map[string]*analysis.Analyzer)

	if marker.EnableAll {
		for _, analyzer := range es.m.GetAllSupportedMarkers() {
			resultMarkersSet[analyzer.Name()] = analyzer
		}

		return resultMarkersSet
	}

	if marker.DisableAll {
		return make(map[string]*analysis.Analyzer)
	}

	// Load default
	for _, analyzer := range es.m.GetAllSupportedMarkers() {
		if es.m.EnabledByDefault(analyzer.Name()) {
			resultMarkersSet[analyzer.Name()] = analyzer
		}
	}

	for _, name := range marker.Enable {
		analyzer := es.m.GetMarkerAnalyzer(name)
		if analyzer != nil {
			resultMarkersSet[name] = analyzer
		}
	}

	for _, m := range marker.Disable {
		analyzer := es.m.GetMarkerAnalyzer(m)
		if analyzer != nil {
			delete(resultMarkersSet, m)
		}
	}

	return resultMarkersSet
}

func (es EnabledSet) GetEnabledMarkersMap() map[string]*analysis.Analyzer {
	enabledMarkers := es.build(&es.opt.Markers)

	return enabledMarkers
}
