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

	rules.NewCherryPick(&rule.Weights{
		GitHubFlow: rule.NewWeight(0.0204560338060376),
		GitFlow:    rule.NewWeight(0.172141890128531),
		GitLabFlow: rule.NewWeight(0.0546810891458411),
		OneFlow:    rule.NewWeight(0.0534825582657873),
		TrunkBased: rule.NewWeight(0.021875),
	}),
	rules.NewCherryPickRelease(&rule.Weights{
		GitHubFlow: rule.NewWeight(0),
		GitFlow:    rule.NewWeight(0.0833333333333333),
		GitLabFlow: rule.NewWeight(0),
		OneFlow:    rule.NewWeight(0),
		TrunkBased: rule.NewWeight(0),
	}),
	rules.NewCrissCrossMerged(&rule.Weights{
		GitHubFlow: rule.NewWeight(0),
		GitFlow:    rule.NewWeight(0),
		GitLabFlow: rule.NewWeight(0),
		OneFlow:    rule.NewWeight(0),
		TrunkBased: rule.NewWeight(0),
	}),
	rules.NewFeatureBranching(&rule.Weights{
		GitHubFlow: rule.NewWeight(0.0934959349593496),
		GitFlow:    rule.NewWeight(0.344935731059201),
		GitLabFlow: rule.NewWeight(0.4547008547008540),
		OneFlow:    rule.NewWeight(0.380150605629752),
		TrunkBased: rule.NewWeight(0.166666666666667),
	}),
	rules.NewHotfix(&rule.Weights{
		GitHubFlow: rule.NewWeight(0),
		GitFlow:    rule.NewWeight(0),
		GitLabFlow: rule.NewWeight(0.04166666667),
		OneFlow:    rule.NewWeight(0.03062200957),
		TrunkBased: rule.NewWeight(0),
	}),
	rules.NewUnresolved(&rule.Weights{
		GitHubFlow: rule.NewWeight(0),
		GitFlow:    rule.NewWeight(0),
		GitLabFlow: rule.NewWeight(0),
		OneFlow:    rule.NewWeight(0),
		TrunkBased: rule.NewWeight(0),
	}),

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
