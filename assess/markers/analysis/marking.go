package analysis

import (
	"github.com/Git-Gopher/go-gopher/assess/options"
	"github.com/Git-Gopher/go-gopher/detector"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/Git-Gopher/go-gopher/violation"
)

type Mark struct {
	Username   string
	Violations []violation.Violation
	Total      int
	Grade      Grade
}

type Grade struct {
	Grade int
	Total int
}

type MarkerCtx struct {
	Model        *enriched.EnrichedModel
	Contribution *contribution
	Author       utils.Authors
}

type MarkerRun func(MarkerCtx) (string, []Mark)

// DetectorMarker is a helper method to run a detector as a marker.
func DetectorMarker(
	m MarkerCtx, // ctx
	ds []detector.Detector, // detector
	c map[string]int, // contribution map
	x int, // contribution multipler
	g options.GradingAlgorithm, // grading algorithm.
) []Mark {
	// violations map to author
	violationsMap := make(map[string][]violation.Violation)
	for _, d := range ds {
		if err := d.Run(m.Model); err != nil {
			return nil
		}

		_, _, _, violations := d.Result()
		for _, v := range violations {
			username, err := m.Author.Find(v.Email())
			if err != nil {
				continue
			}

			violationsMap[*username] = append(violationsMap[*username], v)
		}
	}

	marks := make([]Mark, 0, len(violationsMap))
	for username, violations := range violationsMap {
		emails, err := m.Author.Details(username)
		if err != nil {
			continue
		}

		count := 0 // total contributions for this author.
		for _, email := range emails {
			count += c[email]
		}

		marks = append(marks, Mark{
			Username:   username,
			Violations: violations,
			Total:      count * x,
			Grade: Grade{
				Grade: g(len(violations), count*x),
				Total: 3,
			},
		})
	}

	return marks
}