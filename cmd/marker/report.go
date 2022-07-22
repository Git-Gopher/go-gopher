package main

import (
	"fmt"
	"os"

	"github.com/Git-Gopher/go-gopher/assess"
	"github.com/Git-Gopher/go-gopher/markup"
	"github.com/gomarkdown/markdown"
)

func IndividualReports(candidates []assess.Candidate) error {
	if len(candidates) == 0 {
		return fmt.Errorf("no candidates") // nolint: goerr113
	}

	markers := make([]string, len(candidates[0].Grades))
	for i, grade := range candidates[0].Grades {
		markers[i] = grade.Name
	}

	for _, candidate := range candidates {
		course := "TEST"

		rows := make([][]string, len(candidate.Grades))
		for i, grade := range candidate.Grades {
			rows[i] = []string{
				grade.Name,
				fmt.Sprintf("%d", grade.Violation),
				fmt.Sprintf("%d", grade.Contribution),
				fmt.Sprintf("%d/3", grade.Grade),
			}
		}

		headRow := [][]string{{"Marker", "Violation", "Contribution", "Grade"}}

		rows = append(headRow, rows...)

		m := markup.NewMarkdown()
		md := m.
			Title(fmt.Sprintf("%s %s Report", course, candidate.Username)).
			SubTitle("Marked by git-gopher").
			Table(rows...)

		for _, grade := range candidate.Grades {
			subcategory := fmt.Sprintf("Marker details %s", grade.Name)
			md.SubSubTitle(subcategory).Text(grade.Details)
		}

		output := markdown.ToHTML([]byte(md.String()), nil, nil)

		filename := fmt.Sprintf("%s-individual-reports.html", candidate.Username)

		if err := writeFile(filename, output); err != nil {
			return fmt.Errorf("failed to write file for %s: %w", filename, err)
		}
	}

	return nil
}

func writeFile(filename string, data []byte) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}
