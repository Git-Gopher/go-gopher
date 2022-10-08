package rules

import (
	"fmt"
	"math"

	"github.com/Git-Gopher/go-gopher/detector"
	"github.com/Git-Gopher/go-gopher/workflow/rules/rule"
)

const featureBranchingName = "featureBranching"

func NewFeatureBranching(w *rule.Weights) *rule.Runner {
	runner := rule.NewRunner(
		featureBranchingName,
		"Feature Branching rule",
		func(ctx rule.RuleCtx) (string, *rule.Scores) {
			d := detector.NewFeatureBranchDetector("Feature Branch Detector")

			score := 0.0

			if err := func() error { //nolint:staticcheck
				if err := d.Run(ctx.Model); err != nil {
					return fmt.Errorf("could not run detector: %w", err)
				}

				_, found, total, _ := d.Result()

				score = float64(found) / float64(total)
				if math.IsNaN(score) {
					score = 0
				}

				return nil
			}(); err != nil {
				// do nothing
			}

			return "Feature Branching", w.NewScores(rule.NewScore(score))
		},
	)

	return runner
}
