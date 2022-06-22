package detector

import (
	"testing"

	"github.com/Git-Gopher/go-gopher/model/local"
)

func TestStaleBranchDetect(t *testing.T) {
	r := fetchRepository(t, "https://github.com/Git-Gopher/tests", "test/stale-branch/0")
	tests := []struct {
		name string
		want CommitDetect
	}{
		{"Stale branch", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the model
			model, err := local.NewGitModel(r)
			if err != nil {
				t.Errorf("TestTwoParentsCommitDetect() create model = %v", err)
			}

			detector := NewBranchDetector(StaleBranchDetect())
			if err = detector.Run(model); err != nil {
				t.Errorf("TestStaleBranchDetect() run detector = %v", err)
			}

			t.Log(detector.Result())
		})
	}
}
