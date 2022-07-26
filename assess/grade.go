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

// D Default marker without grading.
func D(m Marker) struct {
	Marker
	*GradingAlgorithm
} {
	return struct {
		Marker
		*GradingAlgorithm
	}{m, nil}
}

// C Custom marker with custom grading algorithm.
func C(m Marker, g GradingAlgorithm) struct {
	Marker
	GradingAlgorithm
} {
	return struct {
		Marker
		GradingAlgorithm
	}{m, g}
}

// nolint: gocognit
func RunMarker(m MarkerCtx, def GradingAlgorithm, markers ...struct {
	Marker
	*GradingAlgorithm
}) []Candidate {
	candiateMap := make(map[string]*Candidate)
	contributionMap := make(map[string]int)

	for _, marker := range markers {
		name, grades := marker.Marker(m)

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

			if _, ok := contributionMap[username]; !ok {
				contributionMap[username] = 0
			}
			contributionMap[username] += contribution
		}

		for username := range candiateMap {
			contribution := contributionMap[username]

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

			details := ""
			for _, v := range grade.Violations {
				details += v.Display(m.Author) + "\n"
			}

			points := 0
			if marker.GradingAlgorithm != nil {
				// use custom grader
				grader := *marker.GradingAlgorithm
				grader(len(grade.Violations), grade.Total)
			} else {
				// default grading
				points = def(len(grade.Violations), grade.Total)
			}

			candiateMap[username].Grades = append(
				candiateMap[username].Grades,
				Grade{
					Name:         name,
					Grade:        points,
					Contribution: grade.Total,
					Violation:    len(grade.Violations),
					Details:      details,
				},
			)

			candiateMap[username].Total += points
		}
	}

	candiates := make([]Candidate, 0, len(candiateMap))
	for _, candidate := range candiateMap {
		candiates = append(candiates, *candidate)
	}

	return candiates
}
