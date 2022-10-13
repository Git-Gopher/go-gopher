package assessors

import (
	"github.com/Git-Gopher/go-gopher/assess/markers"
	"github.com/Git-Gopher/go-gopher/assess/markers/analysis"
	"github.com/Git-Gopher/go-gopher/assess/options"
	log "github.com/sirupsen/logrus"
)

type Manager struct {
	nameToAnalyzer map[string]*analysis.Analyzer
	opt            *options.Options
}

func NewManager(log *log.Logger, opt *options.Options) *Manager {
	m := &Manager{
		nameToAnalyzer: make(map[string]*analysis.Analyzer),
		opt:            opt,
	}

	for _, analyzer := range m.GetAllSupportedMarkers() {
		m.nameToAnalyzer[analyzer.Name()] = analyzer
	}

	return m
}

func (m Manager) GetMarkerAnalyzer(name string) *analysis.Analyzer {
	if analyzer, ok := m.nameToAnalyzer[name]; ok {
		return analyzer
	}

	return nil
}

func (m Manager) GetAllSupportedMarkers() []*analysis.Analyzer {
	var (
		commitOpt      *options.CommitSettings
		branchingOpt   *options.BranchingSettings
		pullRequestOpt *options.PullRequestSettings
		generalOpt     *options.GeneralSettings
	)

	if m.opt != nil {
		commitOpt = &m.opt.MarkersSettings.Commit
		branchingOpt = &m.opt.MarkersSettings.Branching
		pullRequestOpt = &m.opt.MarkersSettings.PullRequest
		generalOpt = &m.opt.MarkersSettings.General
	}

	analyzers := []*analysis.Analyzer{
		markers.NewCommit(commitOpt),
		markers.NewBranching(branchingOpt),
		markers.NewPullRequest(pullRequestOpt),
		markers.NewGeneral(generalOpt),
	}

	return analyzers
}

func (m Manager) EnabledByDefault(name string) bool {
	enabledByDefault := map[string]bool{
		markers.NewCommit(nil).Name():      true,
		markers.NewBranching(nil).Name():   true,
		markers.NewPullRequest(nil).Name(): true,
		markers.NewGeneral(nil).Name():     true,
	}

	return enabledByDefault[name]
}
