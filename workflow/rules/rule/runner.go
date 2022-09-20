package rule

// Weight - weight of the rule, can be negative or positive.
type Weight float64

func (w *Weight) Value() float64 {
	return float64(*w)
}

func NewWeight(weight float64) *Weight {
	w := Weight(weight)

	return &w
}

// Weights - weights for each workflow.
type Weights struct {
	GitHubFlow *Weight
	GitFlow    *Weight
	GitLabFlow *Weight
	OneFlow    *Weight
	TrunkBased *Weight
}

func (w *Weights) NewScores(v *Score) *Scores {
	return &Scores{
		gitHubFlow: v.Weight(w.GitHubFlow),
		gitFlow:    v.Weight(w.GitFlow),
		gitLabFlow: v.Weight(w.GitLabFlow),
		oneFlow:    v.Weight(w.OneFlow),
		trunkBased: v.Weight(w.TrunkBased),
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
	gitHubFlow *Score
	gitFlow    *Score
	gitLabFlow *Score
	oneFlow    *Score
	trunkBased *Score
}

func (s *Scores) GitHubFlow() *Score {
	return s.gitHubFlow
}

func (s *Scores) GitFlow() *Score {
	return s.gitFlow
}

func (s *Scores) GitLabFlow() *Score {
	return s.gitLabFlow
}

func (s *Scores) OneFlow() *Score {
	return s.oneFlow
}

func (s *Scores) TrunkBased() *Score {
	return s.trunkBased
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
