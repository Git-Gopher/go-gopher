package local

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestFetchChunk(t *testing.T) {
	type args struct {
		from *object.Commit
		to   *object.Commit
	}
	tests := []struct {
		name    string
		args    args
		want    []Chunk
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FetchDiffs(tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchChunk() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchChunk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExampleFetchChunk(t *testing.T) {
	utils.Environment("../../.env")
	r := utils.FetchRepository(t, "https://github.com/Git-Gopher/tests", "test/stale-branch/0")

	cIter, err := r.CommitObjects()
	if err != nil {
		t.Errorf("CommitObjects() error = %v", err)
	}

	err = cIter.ForEach(func(c *object.Commit) error {
		if c == nil {
			return fmt.Errorf("CommitObject()")
		}

		for _, ph := range c.ParentHashes {
			// Last parent in the parent hashes is considered to be the 'only' parent
			parent, err := r.CommitObject(ph)
			if err != nil {
				return fmt.Errorf("CommitObject() = %w", err)
			}

			diffs, err := FetchDiffs(parent, c)
			if err != nil {
				return fmt.Errorf("FetchDiffs() = %w", err)
			}

			for _, diff := range diffs {
				fmt.Printf("Name: %s\n", diff.Name)
				fmt.Printf("Added: %s\n", diff.Addition)
				fmt.Printf("Deleted: %s\n", diff.Deletion)
				fmt.Printf("Equal: %s\n", diff.Equal)
			}

		}

		return nil
	})
	if err != nil {
		t.Errorf("cIter() = %v", err)
	}

	t.Fail()
}
