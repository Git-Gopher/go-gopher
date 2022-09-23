package rules

import (
	"fmt"

	"github.com/Git-Gopher/go-gopher/detector"
	"github.com/Git-Gopher/go-gopher/workflow/rules/rule"
)

const unresolvedName = "unresolved"

func NewUnresolved(w *rule.Weights) *rule.Runner {
	runner := rule.NewRunner(
		unresolvedName,
		"Unresolved rule",
		func(ctx rule.RuleCtx) (string, *rule.Scores) {
			d := detector.NewCommitDetector(detector.UnresolvedDetect())

			score := 0.0

			if err := func() error { //nolint:staticcheck
				if err := d.Run(ctx.Model); err != nil {
					return fmt.Errorf("could not run detector: %w", err)
				}

				violated, _, total, _ := d.Result()

				score = float64(violated) / float64(total)

				return nil
			}(); err != nil {
				// do nothing
			}

			return "Unresolved", w.NewScores(rule.NewScore(score))
		},
	)

	return runner
}
