package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Git-Gopher/go-gopher/assess"
	"github.com/Git-Gopher/go-gopher/assess/options"
	"github.com/Git-Gopher/go-gopher/markup"
	"github.com/gomarkdown/markdown"
)

var (
	ErrNoCandidates = errors.New("no candidates")
	ErrCSVWrite     = errors.New("failed to write to csv")
)

const (
	course = "SOFTENG206"
	path   = "marker-report.csv"
)

func IndividualReports(options *options.Options, repoName string, candidates []assess.Candidate) error {
	if len(candidates) == 0 {
		return fmt.Errorf("no candidates") //nolint: goerr113
	}

	markers := make([]string, len(candidates[0].Grades))
	for i, grade := range candidates[0].Grades {
		markers[i] = grade.Name
	}

	for _, candidate := range candidates {
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

		md := markup.CreateMarkdown(fillTemplate(options.HeaderTemplate, candidate.Username, repoName)).
			Header("Marked by git-gopher", 2).
			Table(header, rows)

		for _, grade := range candidate.Grades {
			subcategory := fmt.Sprintf("Marker details %s", grade.Name)
			md.Header(subcategory, 2).Paragraph(grade.Details)
		}

		output := markdown.ToHTML([]byte(md.Render()), nil, nil)

		filename := fillTemplate(options.HeaderTemplate, candidate.Username, repoName) + ".html"
		if len(options.OutputDir) != 0 {
			if _, err := os.Stat(options.OutputDir); errors.Is(err, os.ErrNotExist) {
				if err2 := os.MkdirAll(options.OutputDir, os.ModePerm); err2 != nil {
					return fmt.Errorf("can't create options dir: %w", err)
				}
			}
			filename = filepath.Join(options.OutputDir, filename)
		}

		if err := writeFile(filename, output); err != nil {
			return fmt.Errorf("failed to write file for %s: %w", filename, err)
		}
	}

	return nil
}

// Report for the marker so that they can glean an overview of grading for the course.
// Can then adjust marker as needed.
func MarkerReport(candidates []assess.Candidate) error {
	// Average grade per grade category and overall average grade across categories.
	averageGrade := make(map[string]float64)
	noGrades := make(map[string]int)
	for _, c := range candidates {
		for _, g := range c.Grades {
			averageGrade[g.Name] += float64(g.Grade)
			noGrades[g.Name]++
		}
	}

	overallAverageGrade := 0.0
	for k, v := range averageGrade {
		averageGrade[k] = v / float64(noGrades[k])
		overallAverageGrade += averageGrade[k]
	}

	// Write a report for the marker to CSV.
	var fh *os.File
	// nolint: gosec
	fh, err := os.OpenFile(filepath.Clean(path), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("failed to create marker summary csv: %w", err)
	}
	// nolint: errcheck, gosec
	defer fh.Close()

	// Average Grade table per grade category.
	w := csv.NewWriter(fh)
	var gradeTable [][]string // nolint: prealloc
	gradeTable = append(gradeTable, []string{
		"Grade Category",
		"Average Grade (/3)",
		"Average Grade (%)",
	})
	for k, v := range averageGrade {
		gradeTable = append(gradeTable,
			[]string{
				k,
				fmt.Sprintf("%2.2f", v),
				fmt.Sprintf("%2.2f", (v/3.0)*100),
			})
	}

	if err = w.WriteAll(gradeTable); err != nil {
		return ErrCSVWrite
	}

	// Whitespace.
	if err = w.Write([]string{}); err != nil {
		return ErrCSVWrite
	}

	// Overall average grade table
	if err = w.Write([]string{
		fmt.Sprintf("Overall Average Grade (/%d)", len(averageGrade)*3),
		"Overall Average Grade (%)",
	}); err != nil {
		return ErrCSVWrite
	}
	if err = w.Write([]string{
		fmt.Sprintf("%2.2f", overallAverageGrade),
		fmt.Sprintf("%2.2f", (overallAverageGrade/(float64(len(averageGrade)*3)))*100),
	}); err != nil {
		return ErrCSVWrite
	}

	// Whitespace.
	if err = w.Write([]string{}); err != nil {
		return ErrCSVWrite
	}

	// Candidate table.
	var candidateTable [][]string //nolint: prealloc
	candidateTable = append(candidateTable,
		[]string{
			"Username",
			fmt.Sprintf("OverallGrade (/%d)", len(averageGrade)*3),
			"OverallGrade (%)",
		})

	for _, c := range candidates {
		var row []string
		row = append(row, c.Username)

		var overallGrade float64
		for _, g := range c.Grades {
			overallGrade += float64(g.Grade)
		}

		row = append(row, fmt.Sprintf("%2.2f", overallGrade))
		row = append(row, fmt.Sprintf("%2.2f", (overallGrade/(float64(len(averageGrade)*3)))*100))
		candidateTable = append(candidateTable, row)
	}

	if err = w.WriteAll(candidateTable); err != nil {
		return ErrCSVWrite
	}
	w.Flush()

	return nil
}

func writeFile(path string, data []byte) (err error) {
	f, err := os.Create(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer func() {
		err = f.Close()
	}()

	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", path, err)
	}

	return nil
}

func fillTemplate(template string, username string, repository string) string {
	template = strings.ReplaceAll(template, "{{.Username}}", username)
	template = strings.ReplaceAll(template, "{{ .Username }}", username)

	template = strings.ReplaceAll(template, "{{.Repository}}", repository)
	template = strings.ReplaceAll(template, "{{ .Repository }}", repository)

	return template
}

