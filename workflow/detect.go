package workflow

import (
	"fmt"

	"github.com/Git-Gopher/go-gopher/workflow/rules"
	"github.com/Git-Gopher/go-gopher/workflow/rules/rule"
)

// rulesConfig is set manually.
var rulesConfig = []*rule.Runner{
	// rules.NewExample(&rule.Weights{
	// 	GitHubFlow: rule.NewWeight(1.0),
	// 	GitFlow:    rule.NewWeight(-1.0), // can have negative weights
	// 	GitLabFlow: rule.NewWeight(0.0),  // can have zero weights
	// 	OneFlow:    nil,                  // can be nil, which means that the rule is not used
	// 	TrunkBased: nil,
	// }),
	// rules.NewFeatureBranching(&rule.Weights{
	// 	GitHubFlow: rule.NewWeight(1.0),
	// 	GitFlow:    rule.NewWeight(1.0),
	// 	GitLabFlow: rule.NewWeight(1.0),
	// 	OneFlow:    rule.NewWeight(1.0),
	// 	TrunkBased: rule.NewWeight(-1.0), // feature branching is not used in trunk-based
	// }),
	// rules.NewCherryPick(&rule.Weights{
	// 	GitHubFlow: nil,
	// 	GitFlow:    rule.NewWeight(1.0),
	// 	GitLabFlow: rule.NewWeight(2.0), // Cherry pick is more important in GitLab
	// 	OneFlow:    rule.NewWeight(1.0),
	// 	TrunkBased: nil,
	// }),
	// rules.NewCherryPickRelease(&rule.Weights{
	// 	GitHubFlow: rule.NewWeight(-1.0),
	// 	GitFlow:    rule.NewWeight(1.0),
	// 	GitLabFlow: rule.NewWeight(1.0),
	// 	OneFlow:    rule.NewWeight(1.0),
	// 	TrunkBased: nil,
	// }),
	// rules.NewHotfix(&rule.Weights{
	// 	GitHubFlow: nil,
	// 	GitFlow:    rule.NewWeight(1.0),
	// 	GitLabFlow: rule.NewWeight(1.0),
	// 	OneFlow:    rule.NewWeight(1.0),
	// 	TrunkBased: nil,
	// }),

	rules.NewCherryPickRelease(rule.NewDefaultWeights()),
	rules.NewCherryPick(rule.NewDefaultWeights()),
	rules.NewCrissCrossMerged(rule.NewDefaultWeights()),
	rules.NewFeatureBranching(rule.NewDefaultWeights()),
	rules.NewHotfix(rule.NewDefaultWeights()),
	rules.NewUnresolved(rule.NewDefaultWeights()),

	// When just starting out use default weight before calibrating
	// default weight is 1.0 for all workflows
	// ```go
	// rules.NewExample(rule.NewDefaultWeights()),
	// ```
}

func Detect(ctx rule.RuleCtx) map[string]*rule.Scores {
	results := make(map[string]*rule.Scores)

	for _, r := range rulesConfig {
		name, scores := r.Run(ctx)

		i := 2
		n := name
		for {
			if _, ok := results[n]; !ok {
				break
			}
			n = fmt.Sprintf("%s(%d)", name, i)
			i += 1
		}

		results[n] = scores
	}

	return results
}
