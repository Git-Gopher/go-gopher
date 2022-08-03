package local

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var (
	splitLinesRegexp = regexp.MustCompile(`[^\n]*(\n|$)`)
)

func FetchDiffs(patch *object.Patch) ([]Diff, error) {
	filePatches := patch.FilePatches()

	diffs := make([]Diff, 0)

	for _, fp := range filePatches {
		chunks := fp.Chunks()

		var name string
		from, to := fp.Files()
		if from == nil {
			// New File is created.
			name = to.Path()
		} else if to == nil {
			// File is deleted.
			name = from.Path()
		} else if from.Path() != to.Path() {
			// File is renamed. Not supported.
			// cs.Name = fmt.Sprintf("%s => %s", from.Path(), to.Path())
		} else {
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

	equal := equalSb.String()
	add := addSb.String()
	delete := deleteSb.String()

	return equal, add, delete, nil
}

func DefragmentToDiffPoint(chunks []diff.Chunk) ([]DiffPoint, error) {
	diffPoints := make([]DiffPoint, len(chunks))

	fromLine := 0
	toLine := 0

	for i, chunk := range chunks {
		lines := splitLines(chunk.Content())
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

func splitLines(s string) []string {
	out := splitLinesRegexp.FindAllString(s, -1)
	if out[len(out)-1] == "" {
		out = out[:len(out)-1]
	}
	return out
}
