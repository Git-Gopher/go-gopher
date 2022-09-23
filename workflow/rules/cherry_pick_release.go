package rules

import (
	"fmt"
	"math"

	"github.com/Git-Gopher/go-gopher/detector"
	"github.com/Git-Gopher/go-gopher/workflow/rules/rule"
)

const cherryPickReleaseName = "cherry-pick-release"

func NewCherryPickRelease(w *rule.Weights) *rule.Runner {
	runner := rule.NewRunner(
		cherryPickReleaseName,
		"Cherry pick rule",
		func(ctx rule.RuleCtx) (string, *rule.Scores) {
			d := detector.NewCherryPickReleaseDetector("CherryPickReleaseDetector")

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

			return "CherryPickRelease", w.NewScores(rule.NewScore(score))
		},
	)

	return runner
}
