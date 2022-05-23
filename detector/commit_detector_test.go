package detector

import (
	"reflect"
	"testing"

	"github.com/go-git/go-git/v5/plumbing/object"
)

func MockCommitDetect() CommitDetect {
	return func(commit *object.Commit) (bool, error) {
		return commit.Message == "true", nil
	}
}

func TestCommitDetector_Run(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		commitDetector := NewCommitDetector(MockCommitDetect())

		fakeCommit := &object.Commit{
			Message: "true",
		}

		err := commitDetector.Run(fakeCommit)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// XXX: Do inline
		violated, count, total := commitDetector.Result()
		if violated != 0 {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 1 {
			t.Fatalf("unexpected error: %v", err)
		}
		if total != 1 {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		commitDetector := NewCommitDetector(MockCommitDetect())

		fakeCommit := &object.Commit{
			Message: "false",
		}

		err := commitDetector.Run(fakeCommit)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// XXX: Do inline
		violated, count, total := commitDetector.Result()
		if violated != 0 {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 0 {
			t.Fatalf("unexpected error: %v", err)
		}
		if total != 1 {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestNewLineCommitDetect(t *testing.T) {
	tests := []struct {
		name    string
		commit  *object.Commit
		wantErr bool
		want    bool
	}{
		{"less 10", &object.Commit{Message: "asdf"}, false, true},
		{"over 10", &object.Commit{Message: "1234567890123"}, false, false},
	}
	newLineLengthCommitDetect := NewLineLengthCommitDetect()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := newLineLengthCommitDetect(tc.commit)
			if (err != nil) != tc.wantErr {
				t.Errorf("NewLineCommitDetect() error = %v, wantErr %v", err, tc.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("NewLineCommitDetect() = %v, want %v", got, tc.want)
			}
		})
	}
}
