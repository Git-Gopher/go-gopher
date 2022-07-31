package analysis

type Analyzer struct {
	name, desc string
	Run        func(m MarkerCtx) (string, []Mark)
}

func NewAnalyzer(name, desc string, run func(m MarkerCtx) (string, []Mark)) *Analyzer {
	return &Analyzer{
		name: name,
		desc: desc,
		Run:  run,
	}
}

func (a *Analyzer) Name() string {
	return a.name
}

func (a *Analyzer) Desc() string {
	return a.desc
}
