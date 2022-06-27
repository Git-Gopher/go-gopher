package local

import (
	"fmt"
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func FetchDiffs(from *object.Commit, to *object.Commit) ([]Diff, error) {
	patch, err := from.Patch(to)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch chunk: %w", err)
	}

	files, _, err := gitdiff.Parse(strings.NewReader(patch.String()))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse diff: %w", err)
	}
	diffs := make([]Diff, len(files))
	for i, f := range files {
		equal, added, deleted, err := Defragment(f.TextFragments)
		if err != nil {
			return nil, fmt.Errorf("Failed to defragment diffs: %w", err)
		}

		diffPoints, err := DefragmentToDiffPoint(f.TextFragments)
		if err != nil {
			return nil, fmt.Errorf("Failed to defragment diff points: %w", err)
		}

		var name string
		if f.IsNew || f.IsRename {
			name = f.NewName
		} else {
			name = f.OldName
		}

		diffs[i] = Diff{
			Name:     name,
			Addition: added,
			Deletion: deleted,
			Equal:    equal,

			Points: diffPoints,
		}
	}

	return diffs, nil
}

func Defragment(fragment []*gitdiff.TextFragment) (equal, added, deleted string, err error) {
	for _, f := range fragment {
		for _, l := range f.Lines {
			switch l.Op {
			case gitdiff.LineOp(Equal):
				equal += l.Line
			case gitdiff.LineOp(Add):
				added += l.Line
			case gitdiff.LineOp(Delete):
				deleted += l.Line
			default:
				return "", "", "", ErrUnknownLineOp
			}
		}
	}

	return equal, added, deleted, nil
}

func DefragmentToDiffPoint(fragments []*gitdiff.TextFragment) ([]DiffPoint, error) {
	diffPoints := make([]DiffPoint, len(fragments))
	for i, f := range fragments {
		point := DiffPoint{
			OldPosition:  f.OldPosition,
			OldLines:     f.OldLines,
			NewPosition:  f.NewPosition,
			NewLines:     f.NewLines,
			LinesAdded:   f.LinesAdded,
			LinesDeleted: f.LinesDeleted,
		}
		diffPoints[i] = point
	}

	return diffPoints, nil
}
