package rule

// Weight - weight of the rule, can be negative or positive.
type Weight float64

type WorkflowType int

const (
	// Supported workflow types.
	// XXX: Ordering matters.
	GitHubFlow WorkflowType = iota
	GitFlow
	GitlabFlow
	OneFlow
	TrunkBased
)

// WorkflowType string lookup.
func (wt WorkflowType) String() string {
	return [...]string{
		"GitHubFlow",
		"GitFlow",
		"GitlabFlow",
		"OneFlow",
		"Trunkbased",
	}[wt]
}

func (w *Weight) Value() float64 {
	return float64(*w)
}

func NewWeight(weight float64) *Weight {
	w := Weight(weight)

	return &w
}

// Weights - weights for each workflow.
type Weights struct {
	GitHubFlow *Weight `json:",omitempty"`
	GitFlow    *Weight `json:",omitempty"`
	GitLabFlow *Weight `json:",omitempty"`
	OneFlow    *Weight `json:",omitempty"`
	TrunkBased *Weight `json:",omitempty"`
}

func (w *Weights) NewScores(v *Score) *Scores {
	return &Scores{
		GitHubFlow: v.Weight(w.GitHubFlow),
		GitFlow:    v.Weight(w.GitFlow),
		GitLabFlow: v.Weight(w.GitLabFlow),
		OneFlow:    v.Weight(w.OneFlow),
		TrunkBased: v.Weight(w.TrunkBased),
	}
}

func NewDefaultWeights() *Weights {
	return &Weights{
		GitHubFlow: NewWeight(1.0),
		GitFlow:    NewWeight(1.0),
		GitLabFlow: NewWeight(1.0),
		OneFlow:    NewWeight(1.0),
		TrunkBased: NewWeight(1.0),
	}
}

// Score - score of the rule.
type Score float64

func (s *Score) Value() float64 {
	return float64(*s)
}

func (s *Score) Weight(w *Weight) *Score {
	if w == nil {
		return NewScore(0.0)
	}

	return NewScore(s.Value() * w.Value())
}

func NewScore(score float64) *Score {
	s := Score(score)

	return &s
}

// Scores - scores for each workflow.
type Scores struct {
	GitHubFlow *Score `json:",omitempty"`
	GitFlow    *Score `json:",omitempty"`
	GitLabFlow *Score `json:",omitempty"`
	OneFlow    *Score `json:",omitempty"`
	TrunkBased *Score `json:",omitempty"`
}

// Runner - rule runner.
type Runner struct {
	name, desc string
	Run        RuleRun
}

func NewRunner(name, desc string, run RuleRun) *Runner {
	return &Runner{
		name: name,
		desc: desc,
		Run:  run,
	}
}

func (a *Runner) Name() string {
	return a.name
}

func (a *Runner) Desc() string {
	return a.desc
}
