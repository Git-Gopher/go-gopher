package commands

import (
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

		filename := fillTemplate(options.FilenameTemplate, candidate.Username, repoName) + ".html"
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

func fillTemplate(template string, username string, repository string) string {
	template = strings.ReplaceAll(template, "{{.Username}}", username)
	template = strings.ReplaceAll(template, "{{ .Username }}", username)

	template = strings.ReplaceAll(template, "{{.Repository}}", repository)
	template = strings.ReplaceAll(template, "{{ .Repository }}", repository)

	return template
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
