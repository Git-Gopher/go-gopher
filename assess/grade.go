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
		for _, grade := range grades {
			if _, ok := candiateMap[grade.Username]; !ok {
				candiateMap[grade.Username] = &Candidate{
					Username: grade.Username,
					Grades:   []Grade{},
					Total:    0,
				}
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
