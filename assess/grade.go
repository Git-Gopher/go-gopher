package assess

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

func RunMarker(m MarkerCtx, g GradingAlgorithm, markers ...Marker) []Candidate {
	candiateMap := make(map[string]*Candidate)

	for _, marker := range markers {
		name, grades := marker(m)

		gradeMap := make(map[string]Mark)
		for _, grade := range grades {
			gradeMap[grade.Username] = grade
		}

		for email, contribution := range m.Contribution.CommitCountMap {
			usernamePointer, err := m.Author.Find(email)
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

			var grade Mark
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

				continue
			}

			details := ""
			for _, v := range grade.Violations {
				details += v.Display(m.Author) + "\n"
			}

			points := g(len(grade.Violations), grade.Total)

			candiateMap[grade.Username].Grades = append(
				candiateMap[grade.Username].Grades,
				Grade{
					Name:         name,
					Grade:        points,
					Contribution: grade.Total,
					Violation:    len(grade.Violations),
					Details:      details,
				},
			)

			candiateMap[grade.Username].Total += points
		}
	}

	candiates := make([]Candidate, 0, len(candiateMap))
	for _, candidate := range candiateMap {
		candiates = append(candiates, *candidate)
	}

	return candiates
}