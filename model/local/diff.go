package local

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/scorpionknifes/go-pcre"
)

var splitLinesRegexp = pcre.MustCompileJIT(`[^\n]*(\n|$)`, 0, pcre.STUDY_JIT_COMPILE)

func FetchDiffs(patch *object.Patch) ([]Diff, error) {
	filePatches := patch.FilePatches()

	diffs := make([]Diff, 0)

	for _, fp := range filePatches {
		chunks := fp.Chunks()

		var name string
		from, to := fp.Files()

		switch {
		case from == nil:
			// New File is created.
			// XXX: This panics sometimes on previous repos.
			// Quick workaround is the hardcode below. You should enable it yourself though.
			// Needs to be investigated.
			//name = ""
			name = to.Path()
		case to == nil:
			// File is deleted.
			name = from.Path()
		case from.Path() != to.Path():
			// File is renamed. Not supported.
			// cs.Name = fmt.Sprintf("%s => %s", from.Path(), to.Path())
		default:
			name = from.Path()
		}
		if len(chunks) == 0 {
			// chunk len == 0 means patch is binary.
			diffs = append(diffs, Diff{
				Name:     name,
				IsBinary: fp.IsBinary(),
			})

			continue
		}

		equal, added, deleted, err := Defragment(chunks)
		if err != nil {
			return nil, fmt.Errorf("failed to defragment: %w", err)
		}

		diffPoints, err := DefragmentToDiffPoint(chunks)
		if err != nil {
			return nil, fmt.Errorf("failed to defragment to diff point: %w", err)
		}

		diffs = append(diffs, Diff{
			Name:     name,
			Addition: added,
			Deletion: deleted,
			Equal:    equal,

			Points: diffPoints,
		})
	}

	return diffs, nil
}

func Defragment(chunks []diff.Chunk) (string, string, string, error) {
	var addSb strings.Builder
	var deleteSb strings.Builder
	var equalSb strings.Builder

	for _, chunk := range chunks {
		s := chunk.Content()
		if len(s) == 0 {
			continue
		}

		switch chunk.Type() {
		case diff.Add:
			addSb.WriteString(s)
			addSb.WriteString("\n")
		case diff.Delete:
			deleteSb.WriteString(s)
			deleteSb.WriteString("\n")
		case diff.Equal:
			equalSb.WriteString(s)
			equalSb.WriteString("\n")
		}
	}

	return equalSb.String(), addSb.String(), deleteSb.String(), nil
}

func DefragmentToDiffPoint(chunks []diff.Chunk) ([]DiffPoint, error) {
	diffPoints := make([]DiffPoint, len(chunks))

	fromLine := 0
	toLine := 0

	for i, chunk := range chunks {
		lines, err := splitLines(chunk.Content())
		if err != nil {
			return nil, fmt.Errorf("failed to split lines: %w", err)
		}
		nLines := len(lines)

		s := chunk.Content()
		if len(s) == 0 {
			diffPoints[i] = DiffPoint{}

			continue
		}

		linesAdded := 0
		linesDeleted := 0

		switch chunk.Type() {
		case diff.Add:
			fromLine += nLines
			toLine += nLines
			linesAdded = nLines
		case diff.Delete:
			if nLines != 0 {
				fromLine++
			}
			fromLine += nLines - 1
			linesDeleted = nLines
		case diff.Equal:
			if nLines != 0 {
				toLine++
			}
			toLine += nLines - 1
		}

		point := DiffPoint{
			OldPosition:  int64(fromLine),
			NewPosition:  int64(toLine),
			LinesAdded:   int64(linesAdded),
			LinesDeleted: int64(linesDeleted),
		}
		diffPoints[i] = point
	}

	return diffPoints, nil
}

func splitLines(s string) ([]string, error) {
	matches, err := splitLinesRegexp.FindAll(s, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to find all: %w", err)
	}

	if matches[len(matches)-1].Finding == "" {
		matches = matches[:len(matches)-1]
	}

	lines := make([]string, len(matches))
	for i, m := range matches {
		lines[i] = m.Finding
	}

	return lines, nil
}
