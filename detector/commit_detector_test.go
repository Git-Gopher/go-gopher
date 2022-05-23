package detector

import (
	"reflect"
	"testing"
)

func MockCommitDetect() CommitDetect {
	return func(commit *Commit) (bool, error) {
		return commit.Message == "true", nil
	}
}

func TestCommitDetector_Run(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		commitDetector := NewCommitDetector(MockCommitDetect())

		fakeCommit := &Commit{
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

		fakeCommit := &Commit{
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
		commit  *Commit
		wantErr bool
		want    bool
	}{
<<<<<<< HEAD
		{"less 10", &Commit{Message: "asdf"}, false, true},
		{"over 10", &Commit{Message: "1234567890123"}, false, false},
=======
		{"less 10", &object.Commit{Message: "asdf"}, false, true},
		{"over 10", &object.Commit{Message: "1234567890123"}, false, false},
>>>>>>> db72aa3e1feb3882d08996bd45e9e93304101451
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
