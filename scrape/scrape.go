package scrape

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var basic struct {
	Viewer struct {
		Login     githubv4.String
		AvatarUrl githubv4.String
	}
}

func Scrape() {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	// basic query
	if err := client.Query(context.Background(), &basic, nil); err != nil {
		log.Fatal("Failed to make basic query")

		return
	}
	fmt.Println(basic)
}
