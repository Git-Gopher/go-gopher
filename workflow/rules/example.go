package rules

import (
	"fmt"
	"math"

	"github.com/Git-Gopher/go-gopher/detector"
	"github.com/Git-Gopher/go-gopher/workflow/rules/rule"
)

const exampleName = "example"

func NewExample(w *rule.Weights) *rule.Runner {
	runner := rule.NewRunner(
		exampleName,
		"Example rule",
		func(ctx rule.RuleCtx) (string, *rule.Scores) {
			d := detector.NewCommitDetector(detector.EmptyCommitDetect())

			score := 0.0

			if err := func() error { //nolint:staticcheck
				if err := d.Run(ctx.Model); err != nil {
					return fmt.Errorf("could not run detector: %w", err)
				}

				_, found, total, _ := d.Result()

<<<<<<< HEAD
				if total != 0 {
					score = float64(found) / float64(total)
=======
				score = float64(found) / float64(total)
				if math.IsNaN(score) {
					score = 0
>>>>>>> main
				}

				return nil
			}(); err != nil {
				// do nothing
			}

			return "Example", w.NewScores(rule.NewScore(score))
		},
	)

	return runner
}
