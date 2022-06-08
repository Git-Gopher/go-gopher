package scrape

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/Git-Gopher/go-gopher/model/github"
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
type UserQuery struct {
	Viewer struct {
		Login     githubv4.String
		AvatarUrl githubv4.String
		Email     githubv4.String
	}
}

func (s *Scraper) ScrapeUsers() (*UserQuery, error) {
	userQuery := new(UserQuery)
	if err := s.Client.Query(context.Background(), userQuery, nil); err != nil {
		return nil, ErrUserQuery
	}

	return userQuery, nil
}

// TODO: Pagnation: fetch entire history. Also handle nested pagination
type PullRequestQuery struct {
	Repository struct {
		PullRequests struct {
			Nodes []struct {
				Title                   string
				Body                    string
				ClosingIssuesReferences struct {
					Edges []struct {
						Node struct {
							Id string
						}
					}
				} `graphql:"closingIssuesReferences(first: $first, after: $after)"`
			}
			PageInfo struct {
				StartCursor string
				HasNextPage bool
				EndCursor   string
			}
		} `graphql:"pullRequests(first: 100)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func (s *Scraper) ScrapePullRequests(owner, name string) ([]*github.PullRequest, error) {
	// Limited to maximum of 100
	var first = 100
	var after = ""
	var hasNextPage = true
	var pullRequests []*github.PullRequest

	for hasNextPage {

		pullRequestQuery := new(PullRequestQuery)
		variables := map[string]interface{}{
			"first": githubv4.Int(first),
			"after": githubv4.String(after),
			"owner": githubv4.String(owner),
			"name":  githubv4.String(name),
		}

		if err := s.Client.Query(context.Background(), pullRequestQuery, variables); err != nil {
			return nil, ErrPullRequestQuery
		}

		// Convert to model
		for _, prq := range pullRequestQuery.Repository.PullRequests.Nodes {
			var issues []*github.Issue
			for _, issue := range prq.ClosingIssuesReferences.Edges {
				issues = append(issues, &github.Issue{Id: issue.Node.Id})
			}
			pullRequests = append(pullRequests, &github.PullRequest{Title: prq.Title, Body: prq.Body, Issues: issues})
		}

		if pullRequestQuery.Repository.PullRequests.PageInfo.HasNextPage {
			after = pullRequestQuery.Repository.PullRequests.PageInfo.EndCursor
		} else {
			hasNextPage = false
		}
	}

	return pullRequests, nil
}

// Fetch all records, handling pagnation
// func (s *Scraper) queryAll(ctx context.Context, query interface{}, map[string]interface{}) {
// }
