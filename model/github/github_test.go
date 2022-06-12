package github

import (
	"fmt"
	"testing"

	"github.com/Git-Gopher/go-gopher/utils"
)

// TODO: Move this to test utils?

func TestScrapeGithubModel(t *testing.T) {
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
			model, err := ScrapeGithubModel(tt.args.owner, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScrapeGithubModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("model: %v\n", model)

			// &GithubModel{Author: nil, PullRequests: []PullRequest{{Title: "test/linked-pull-request-issue/modify", Body: "closes #1", Issues: []Issue{{Id: "l", Title: "", Body: "", Author: nil}}}}}, false,

			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("ScrapeGithubModel() = %v, want %v", got, tt.want)
			// }
		})
	}
}
