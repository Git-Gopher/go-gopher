package rules

import (
	"fmt"
	"math"

	"github.com/Git-Gopher/go-gopher/detector"
	"github.com/Git-Gopher/go-gopher/workflow/rules/rule"
)

const emptyCommitName = "empty-commit"

func NewEmptyCommit(w *rule.Weights) *rule.Runner {
	runner := rule.NewRunner(
		emptyCommitName,
		"EmptyCommit rule",
		func(ctx rule.RuleCtx) (string, *rule.Scores) {
			d := detector.NewCommitDetector(detector.EmptyCommitDetect())

			score := 0.0

			if err := func() error { //nolint:staticcheck
				if err := d.Run(ctx.Model); err != nil {
					return fmt.Errorf("could not run detector: %w", err)
				}

				violated, _, total, _ := d.Result()

				score = float64(violated) / float64(total)
				if math.IsNaN(score) {
					score = 0
				}

				return nil
			}(); err != nil {
				// do nothing
			}

			return "EmptyCommit", w.NewScores(rule.NewScore(score))
		},
	)

	return runner
}
