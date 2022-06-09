package github

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// TODO: Move this to test utils?
func environment() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Error loading .env file")
	}

	if os.Getenv("GITHUB_TOKEN") == "" {
		log.Fatalln("Error loading env GITHUB_TOKEN")
	}
}

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
			environment()
			model, err := ScrapeGithubModel(tt.args.owner, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScrapeGithubModel() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Printf("model: %v\n", model)

			// &GithubModel{Author: nil, PullRequests: []PullRequest{{Title: "test/linked-pull-request-issue/modify", Body: "closes #1", Issues: []Issue{{Id: "l", Title: "", Body: "", Author: nil}}}}}, false,

			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("ScrapeGithubModel() = %v, want %v", got, tt.want)
			// }
		})
	}
}
