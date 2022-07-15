package detector

import (
	"testing"

	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/model/remote"
)

func TestPullRequestLinkedIssue(t *testing.T) {
	tests := []struct {
		name string
		want CommitDetect
	}{
		{"TestPullRequestLinkedIssue", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the githubModel
			githubModel, err := remote.ScrapeRemoteModel("Git-Gopher", "tests")
			if err != nil {
				t.Errorf("%s: scrape github model = %v", tt.name, err)
			}

			enrichedModel := enriched.NewEnrichedModel(local.GitModel{}, *githubModel)
			detector := NewPullRequestDetector(PullRequestIssueDetector())
			if err = detector.Run(enrichedModel); err != nil {
				t.Errorf("%s run detector = %v", tt.name, err)
			}

			t.Log(detector.Result())
		})
	}
}
