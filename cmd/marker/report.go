package main

import (
	"fmt"
	"os"
	"path/filepath"

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

		header := []string{"Marker", "Violation", "Contribution", "Grade"}

		md := markup.CreateMarkdown(fmt.Sprintf("%s %s Report", course, candidate.Username)).
			Header("Marked by git-gopher", 2).
			Table(header, rows)

		for _, grade := range candidate.Grades {
			subcategory := fmt.Sprintf("Marker details %s", grade.Name)
			md.Header(subcategory, 2).Paragraph(grade.Details)
		}

		output := markdown.ToHTML([]byte(md.Render()), nil, nil)

		filename := fmt.Sprintf("%s-individual-reports.html", candidate.Username)

		if err := writeFile(filename, output); err != nil {
			return fmt.Errorf("failed to write file for %s: %w", filename, err)
		}
	}

	return nil
}

func writeFile(filename string, data []byte) (err error) {
	f, err := os.Create(filepath.Clean(filename))
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer func() {
		err = f.Close()
	}()

	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filename, err)
	}

	return nil
}
