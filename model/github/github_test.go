package github

import (
	"fmt"
	"testing"
)

func TestScrapeGithubModel(t *testing.T) {
	type args struct {
		remote string
		owner  string
		name   string
	}
	tests := []struct {
		name    string
		args    args
		want    *GithubModel
		wantErr bool
	}{
		{"ScrapeGithubModel", args{remote: "https://github.com/Git-Gopher/tests", owner: "Git-Gopher", name: "tests"},
			&GithubModel{Author: nil, PullRequests: []PullRequest{PullRequest{Title: "test/linked-pull-request-issue/modify", Body: "closes #1", Issues: []Issue{}}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ghm, err := ScrapeGithubModel(tt.args.remote, tt.args.owner, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScrapeGithubModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			fmt.Printf("ghm: %v\n", ghm)

			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("ScrapeGithubModel() = %v, want %v", got, tt.want)
			// }
		})
	}
}
