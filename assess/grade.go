package assess

import (
	"strings"

	"github.com/Git-Gopher/go-gopher/assess/markers/analysis"
)

type Candidate struct {
	Username string  `json:"username"`
	Grades   []Grade `json:"grades"`
	Total    int     `json:"total"`
}

type Grade struct {
	Name         string `json:"name"`
	Grade        int    `json:"grade"`
	Contribution int    `json:"contribution"`
	Violation    int    `json:"violation"`
	Details      string `json:"details"`
}

func RunMarker(m analysis.MarkerCtx, markers []*analysis.Analyzer) []Candidate {
	candiateMap := make(map[string]*Candidate)
	contributionMap := make(map[string]int)

	for _, marker := range markers {
		name, grades := marker.Run(m)

		gradeMap := make(map[string]analysis.Mark)
		for _, grade := range grades {
			gradeMap[grade.Username] = grade
		}

		for email, contribution := range m.Contribution.CommitCountMap {
			usernamePointer, err := m.Author.FindUserName(email)
			if err != nil {
				continue
			}

			username := *usernamePointer

			if _, ok := candiateMap[username]; !ok {
				candiateMap[username] = &Candidate{
					Username: username,
					Grades:   []Grade{},
					Total:    0,
				}
			}

			if _, ok := contributionMap[username]; !ok {
				contributionMap[username] = 0
			}
			contributionMap[username] += contribution
		}

		for username := range candiateMap {
			contribution := contributionMap[username]

			var grade analysis.Mark
			var ok bool
			if grade, ok = gradeMap[username]; !ok {
				// TODO less contribution
				if contribution < 5 {
					candiateMap[username].Grades = append(
						candiateMap[username].Grades,
						Grade{
							Name:         name,
							Grade:        0,
							Contribution: contribution,
							Violation:    0,
							Details:      "Not enough contribution in this category.",
						},
					)

					continue
				}

				candiateMap[username].Grades = append(
					candiateMap[username].Grades,
					Grade{
						Name:         name,
						Grade:        3,
						Contribution: contribution,
						Violation:    0,
						Details:      "",
					},
				)
				candiateMap[username].Total += 3

				continue
			}

			var detailsSb strings.Builder
			for _, v := range grade.Violations {
				detailsSb.WriteString(v.Display(m.Author))
				detailsSb.WriteString("\n")
			}

			candiateMap[username].Grades = append(
				candiateMap[username].Grades,
				Grade{
					Name:         name,
					Grade:        grade.Grade.Grade,
					Contribution: grade.Total,
					Violation:    len(grade.Violations),
					Details:      detailsSb.String(),
				},
			)

			candiateMap[username].Total += grade.Grade.Grade
		}
	}

	candiates := make([]Candidate, 0, len(candiateMap))
	for _, candidate := range candiateMap {
		candiates = append(candiates, *candidate)
	}

	return candiates
}
