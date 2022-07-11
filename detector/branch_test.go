package detector

import (
	"testing"

	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/model/remote"
	"github.com/Git-Gopher/go-gopher/utils"
)

func TestStaleBranchDetect(t *testing.T) {
	r := utils.FetchRepository(t, "https://github.com/Git-Gopher/tests", "test/stale-branch/0")
	tests := []struct {
		name string
		want CommitDetect
	}{
		{"Stale branch", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the gitModel
			gitModel, err := local.NewGitModel(r)
			if err != nil {
				t.Errorf("TestTwoParentsCommitDetect() create model = %v", err)
			}

			detector := NewBranchDetector(StaleBranchDetect())
			enrichedModel := enriched.NewEnrichedModel(*gitModel, remote.RemoteModel{})
			if err = detector.Run(enrichedModel); err != nil {
				t.Errorf("TestStaleBranchDetect() run detector = %v", err)
			}

			t.Log(detector.Result())
		})
	}
}
