package scrape

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var ErrBasicQuery = errors.New("Failed to make basic query")

type Scraper struct {
	Client *githubv4.Client
	Remote string
}

func NewScraper(remote string) Scraper {
	// XXX: extract this out for tests and logic, tests should be using a setup, should be
	// fine for running the application because main entry point sets this up beforehand
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Error loading .env file")
	}

	if os.Getenv("GITHUB_TOKEN") == "" {
		log.Fatalln("Error loading env GITHUB_TOKEN")
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	return Scraper{
		client,
		remote,
	}
}

// TODO: put this inline? Or leave it outside so that we have a nice preview of what we are getting.
type UserQuery struct {
	Viewer struct {
		Login     githubv4.String
		AvatarUrl githubv4.String
	}
}

func (s *Scraper) ScrapeUsers() (*UserQuery, error) {
	userQuery := new(UserQuery)
	if err := s.Client.Query(context.Background(), userQuery, nil); err != nil {
		return nil, ErrBasicQuery
	}

	return userQuery, nil
}
