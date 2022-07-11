package remote

import (
	"testing"

	"github.com/Git-Gopher/go-gopher/utils"
)

func TestScrapeRemoteModel(t *testing.T) {
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
			"ScrapeRemoteModel",
			args{owner: "Git-Gopher", name: "tests"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utils.Environment("../../.env")
			_, err := ScrapeRemoteModel(tt.args.owner, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScrapeRemoteModel() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
		})
	}
}

func BenchmarkScrapeRemoteModel(b *testing.B) {
	utils.Environment("../../.env")
	model, err := ScrapeRemoteModel("subquery", "subql") // Just an example of a large repository
	if err != nil {
		b.Errorf("BenchmarkScrapeRemoteModel() error = %v", err)
	}

	b.Logf("BenchmarkScrapeRemoteModel() Issues = %d", len(model.Issues))
	b.Logf("BenchmarkScrapeRemoteModel() Pull request = %d", len(model.PullRequests))
	comments := make([]*Comment, 0)
	for _, pr := range model.PullRequests {
		comments = append(comments, pr.Comments...)
	}
	for _, is := range model.Issues {
		comments = append(comments, is.Comments...)
	}

	b.Logf("BenchmarkScrapeRemoteModel() Total comments = %d", len(comments))
}
