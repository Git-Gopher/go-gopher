package github

import (
	"testing"

	"github.com/Git-Gopher/go-gopher/utils"
)

func TestScrapeGitHubModel(t *testing.T) {
	type args struct {
		owner string
		name  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"ScrapeGithubModel",
			args{owner: "Git-Gopher", name: "tests"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utils.Environment("../../.env")
			_, err := ScrapeGithubModel(tt.args.owner, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScrapeGithubModel() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
		})
	}
}

func BenchmarkScrapeGitHubModel(b *testing.B) {
	utils.Environment("../../.env")
	model, err := ScrapeGithubModel("subquery", "subql") // Just an example of a large repository
	if err != nil {
		b.Errorf("BenchmarkScrapeGitHubModel() error = %v", err)
	}

	b.Logf("BenchmarkScrapeGitHubModel() Issues = %d", len(model.Issues))
	b.Logf("BenchmarkScrapeGitHubModel() Pull request = %d", len(model.PullRequests))
	comments := make([]*Comment, 0)
	for _, pr := range model.PullRequests {
		comments = append(comments, pr.Comments...)
	}
	for _, is := range model.Issues {
		comments = append(comments, is.Comments...)
	}

	b.Logf("BenchmarkScrapeGitHubModel() Total comments = %d", len(comments))
}
