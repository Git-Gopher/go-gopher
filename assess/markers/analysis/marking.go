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
	Model          *enriched.EnrichedModel
	Contribution   *contribution
	Author         utils.Authors
	LoginWhiteList []string
	Upis           map[string]string
	Fullnames      map[string]string
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

		violations = violation.FilterByLogin(violations, nil, m.LoginWhiteList)

		for _, v := range violations {
			login := ""
			if l, err := v.Login(); err == nil {
				login = l
			} else if email, err := v.Email(); err == nil {
				l, err := m.Author.FindUserName(email)
				if err != nil {
					continue
				}
				login = *l
			} else {
				continue
			}

			violationsMap[login] = append(violationsMap[login], v)
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
