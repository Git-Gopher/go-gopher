package rules

import (
	"fmt"

	"github.com/Git-Gopher/go-gopher/detector"
	"github.com/Git-Gopher/go-gopher/workflow/rules/rule"
)

const cherryPickName = "cherry-pick"

func NewCherryPick(w *rule.Weights) *rule.Runner {
	runner := rule.NewRunner(
		cherryPickName,
		"Cherry pick rule",
		func(ctx rule.RuleCtx) (string, *rule.Scores) {
			d := detector.NewCherryPickDetector("CherryPickDetector")

			score := 0.0

			if err := func() error { //nolint:staticcheck
				if err := d.Run(ctx.Model); err != nil {
					return fmt.Errorf("could not run detector: %w", err)
				}

				_, found, total, _ := d.Result()

				if total != 0 {
					score = float64(found) / float64(total)
				}

				return nil
			}(); err != nil {
				// do nothing
			}

			return "Cherry Pick", w.NewScores(rule.NewScore(score))
		},
	)

	return runner
}
