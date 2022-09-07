package rules

import (
	"fmt"

	"github.com/Git-Gopher/go-gopher/detector"
	"github.com/Git-Gopher/go-gopher/workflow/rules/rule"
)

const taggedMainName = "taggedMain"

func NewTaggedMain(w *rule.Weights) *rule.Runner {
	return rule.NewRunner(
		taggedMainName,
		"Tagged Main Commits",
		func(ctx rule.RuleCtx) (string, *rule.Scores) {
			d := detector.NewCommitDetector(detector.EmptyCommitDetect())

			score := 0.0

			if err := func() error { //nolint:staticcheck
				if err := d.Run(ctx.Model); err != nil {
					return fmt.Errorf("could not run detector: %w", err)
				}

				_, found, total, _ := d.Result()

				score = float64(found) / float64(total)

				return nil
			}(); err != nil {
				// do nothing
			}

			return "Tagged Main Commits", w.NewScores(rule.NewScore(score))
		},
	)
}
