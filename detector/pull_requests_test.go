package detector

import (
	"testing"

	"github.com/Git-Gopher/go-gopher/model/github"
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
			// Create the model
			model, err := github.ScrapeGithubModel("Git-Gopher", "tests")
			if err != nil {
				t.Errorf("%s: scrape github model = %v", tt.name, err)
			}

			detector := NewPullRequestDetector(PullRequestIssueDetector())
			if err = detector.Run2(model); err != nil {
				t.Errorf("%s run detector = %v", tt.name, err)
			}

			t.Log(detector.Result())
		})
	}
}
