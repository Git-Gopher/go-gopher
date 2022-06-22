package github

import (
	"context"
	"fmt"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// dedicated hyrdation functions per use case.
// define queries internally to function, aswell as any needed structs which should also be package only.
// scraper should also be here!

// instead do this:
// Fetch all one query and then fetch more if page info says so!

var GithubQuerySize = 100

type Scraper struct {
	Client *githubv4.Client
}

func NewScraper() Scraper {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	return Scraper{
		client,
	}
}

type PageInfo struct {
	HasNextPage bool
	EndCursor   githubv4.String
}

func (s *Scraper) FetchPullRequestClosingIssues(owner, name string, number int, cursor string) ([]*Issue, error) {
	var q struct {
		Repository struct {
			PullRequest struct {
				ClosingIssuesReferences struct {
					Nodes []struct {
						Id     string
						Title  string
						Body   string
						Author struct {
							Login     string
							AvatarUrl string
						}
					}
					PageInfo PageInfo
				} `graphql:"closingIssuesReferences(first: $first, after: $cursor)"`
			} `graphql:"pullRequest(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	var all []*Issue
	variables := map[string]interface{}{
		"number": githubv4.Int(number),
		"first":  githubv4.Int(GithubQuerySize),
		"cursor": githubv4.String(cursor),
		"owner":  githubv4.String(owner),
		"name":   githubv4.String(name),
	}

	for {
		if err := s.Client.Query(context.Background(), q, variables); err != nil {
			return nil, fmt.Errorf("Failed to fetch additional pull request closing issues references: %v", err)
		}

		for _, i := range q.Repository.PullRequest.ClosingIssuesReferences.Nodes {
			issue := Issue{
				Id:    i.Id,
				Title: i.Title,
				Body:  i.Body,
				Author: &Author{
					Login:     i.Author.Login,
					AvatarUrl: i.Author.AvatarUrl,
				},
			}

			all = append(all, &issue)
		}

		if !q.Repository.PullRequest.ClosingIssuesReferences.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = githubv4.NewString(q.Repository.PullRequest.ClosingIssuesReferences.PageInfo.EndCursor)
	}

	return all, nil
}

func (s *Scraper) FetchPullRequestComments(owner, name string, number int, cursor string) ([]*Comment, error) {
	var q struct {
		Repository struct {
			PullRequest struct {
				Comments struct {
					Nodes []struct {
						Id     string
						Body   string
						Author struct {
							Login     string
							AvatarUrl string
						}
					}
					PageInfo PageInfo
				} `graphql:"comments(first: $first, after: $cursor)"`
			} `graphql:"pullRequest(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	var all []*Comment
	variables := map[string]interface{}{
		"number": githubv4.Int(number),
		"first":  githubv4.Int(GithubQuerySize),
		"cursor": githubv4.String(cursor),
		"owner":  githubv4.String(owner),
		"name":   githubv4.String(name),
	}

	for {
		if err := s.Client.Query(context.Background(), q, variables); err != nil {
			return nil, fmt.Errorf("Failed to fetch additional pull request closing issues references: %v", err)
		}

		for _, i := range q.Repository.PullRequest.Comments.Nodes {
			comment := Comment{
				Id:   i.Id,
				Body: i.Body,
				Author: &Author{
					Login:     i.Author.Login,
					AvatarUrl: i.Author.AvatarUrl,
				},
			}

			all = append(all, &comment)
		}

		if !q.Repository.PullRequest.Comments.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = githubv4.NewString(q.Repository.PullRequest.Comments.PageInfo.EndCursor)
	}

	return all, nil
}

func (s *Scraper) FetchPullRequests(owner, name string) ([]*PullRequest, error) {
	var q struct {
		Repository struct {
			PullRequests struct {
				Nodes []struct {
					// PullRequest
					Id             string
					Number         int
					Title          string
					Body           string
					ReviewDecision string
					Merged         bool
					// Author
					Author struct {
						Login     string
						AvatarUrl string
					}
					// Issues
					ClosingIssuesReferences struct {
						Nodes []struct {
							Id     string
							Title  string
							Body   string
							Author struct {
								Login     string
								AvatarUrl string
							}
						}
						PageInfo PageInfo
					} `graphql:"closingIssuesReferences(first: $first)"`
					// Comments
					Comments struct {
						Nodes []struct {
							Id     string
							Body   string
							Author struct {
								Login     string
								AvatarUrl string
							}
						}
						PageInfo PageInfo
					} `graphql:"comments(first: $first)"`
					// XXX: Limitation of the GraphQL API, can't properly paginate reviewthreads, limiting to 100 for now.
					ReviewThreads struct {
						Nodes []struct {
							Id         string
							IsOutdated bool
							IsResolved bool
							Path       string
						}
					} `graphql:"reviewThreads(first: 100)"`
				}
				PageInfo PageInfo
			} `graphql:"pullRequests(first: $first, after: $cursor)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	var all []*PullRequest
	variables := map[string]interface{}{
		"first":  githubv4.Int(GithubQuerySize),
		"cursor": (*githubv4.String)(nil),
		"owner":  githubv4.String(owner),
		"name":   githubv4.String(name),
	}

	for {
		if err := s.Client.Query(context.Background(), &q, variables); err != nil {
			return nil, fmt.Errorf("Failed to fetch pull requests: %v", err)
		}

		for _, mpr := range q.Repository.PullRequests.Nodes {
			pr := PullRequest{
				Id:             mpr.Id,
				Number:         mpr.Number,
				Title:          mpr.Title,
				Body:           mpr.Body,
				ReviewDecision: mpr.ReviewDecision,
				Merged:         mpr.Merged,
				Author:         (*Author)(&mpr.Author),
				ClosingIssues:  nil,
				Comments:       nil,
				ReviewThreads:  nil,
			}

			// Closing issues
			var cis []*Issue = make([]*Issue, len(mpr.ClosingIssuesReferences.Nodes))
			for i, ci := range mpr.ClosingIssuesReferences.Nodes {
				cis[i] = &Issue{
					Id:     ci.Id,
					Title:  ci.Title,
					Body:   ci.Body,
					Author: (*Author)(&ci.Author),
				}
			}

			if mpr.ClosingIssuesReferences.PageInfo.HasNextPage {
				acir, err := s.FetchPullRequestClosingIssues(owner, name, pr.Number,
					string(mpr.ClosingIssuesReferences.PageInfo.EndCursor))
				if err != nil {
					return nil, fmt.Errorf("Failed to fetch pull request closing issue references: %v", err)
				}

				cis = append(cis, acir...)

			}

			pr.ClosingIssues = cis

			// Comments
			var cs []*Comment = make([]*Comment, len(mpr.Comments.Nodes))
			for i, c := range mpr.Comments.Nodes {
				cs[i] = &Comment{
					Id:     c.Id,
					Body:   c.Body,
					Author: (*Author)(&c.Author),
				}
			}

			if mpr.Comments.PageInfo.HasNextPage {
				acs, err := s.FetchPullRequestComments(owner, name, pr.Number,
					string(mpr.Comments.PageInfo.EndCursor))
				if err != nil {
					return nil, fmt.Errorf("Failed to fetch pull request comments: %v", err)
				}

				cs = append(cs, acs...)
			}

			pr.Comments = cs

			// Review threads
			var rs []*ReviewThread = make([]*ReviewThread, len(mpr.ReviewThreads.Nodes))
			for i, r := range mpr.ReviewThreads.Nodes {
				rt := ReviewThread{
					Id:         r.Id,
					IsResolved: r.IsResolved,
					IsOutdated: r.IsOutdated,
					Path:       r.Path,
				}

				rs[i] = &rt
			}

			pr.ReviewThreads = rs

			all = append(all, &pr)
		}

		if !q.Repository.PullRequests.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = githubv4.NewString(q.Repository.PullRequests.PageInfo.EndCursor)
	}

	return all, nil
}
