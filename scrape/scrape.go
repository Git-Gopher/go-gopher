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

var (
	ErrUserQuery        = errors.New("Failed to make user query")
	ErrPullRequestQuery = errors.New("Failed to make pull request query")
)

type Scraper struct {
	Client *githubv4.Client
	Remote string
}

func NewScraper(remote string) Scraper {
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

type user struct {
	Login     githubv4.String
	AvatarUrl githubv4.String
	Email     githubv4.String
}

type userQuery struct {
	Viewer user
}

func (s *Scraper) ScrapeUsers() (*userQuery, error) {
	userQuery := new(userQuery)
	if err := s.Client.Query(context.Background(), userQuery, nil); err != nil {
		return nil, ErrUserQuery
	}

	return userQuery, nil
}

type issue struct {
	Id    string
	Title string
	Body  string
}

type pullRequest struct {
	Id                      string
	Title                   string
	Body                    string
	ClosingIssuesReferences struct {
		Edges []struct {
			Node issue
		}
	} `graphql:"closingIssuesReferences(first: $first)"`
}

// TODO: Pagination: fetch entire closing issue reference history and handle nested pagination
type pullRequestQuery struct {
	Repository struct {
		PullRequests struct {
			Nodes    []pullRequest
			PageInfo struct {
				HasNextPage bool
				EndCursor   githubv4.String
			}
		} `graphql:"pullRequests(first: $first, after: $cursor)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func (s *Scraper) ScrapePullRequests(owner, name string) ([]pullRequest, error) {
	var allPrs []pullRequest
	variables := map[string]interface{}{
		"first":  githubv4.Int(100), // Limited to maximum 100
		"cursor": (*githubv4.String)(nil),
		"owner":  githubv4.String(owner),
		"name":   githubv4.String(name),
	}

	for {
		prs := new(pullRequestQuery)
		if err := s.Client.Query(context.Background(), prs, variables); err != nil {
			// fmt.Printf("err: %v\n", err)
			return nil, ErrPullRequestQuery
		}

		allPrs = append(allPrs, prs.Repository.PullRequests.Nodes...)
		if !prs.Repository.PullRequests.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(prs.Repository.PullRequests.PageInfo.EndCursor)
	}

	return allPrs, nil
}
