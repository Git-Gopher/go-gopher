package github

import (
	"context"
	"fmt"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var githubQuerySize = 100

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

// nolint
func (s *Scraper) FetchIssueComments(owner, name string, number int, cursor string) ([]*Comment, error) {
	var q struct {
		Repository struct {
			Issue struct {
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
			} `graphql:"issue(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	var all []*Comment
	variables := map[string]interface{}{
		"number": githubv4.Int(number),
		"first":  githubv4.Int(githubQuerySize),
		"cursor": githubv4.String(cursor),
		"owner":  githubv4.String(owner),
		"name":   githubv4.String(name),
	}

	for {
		if err := s.Client.Query(context.Background(), q, variables); err != nil {
			return nil, fmt.Errorf("Failed to fetch additional pull request closing issues references: %w", err)
		}

		for _, i := range q.Repository.Issue.Comments.Nodes {
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

		if !q.Repository.Issue.Comments.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = githubv4.NewString(q.Repository.Issue.Comments.PageInfo.EndCursor)
	}

	return all, nil
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
		"first":  githubv4.Int(githubQuerySize),
		"cursor": githubv4.String(cursor),
		"owner":  githubv4.String(owner),
		"name":   githubv4.String(name),
	}

	for {
		if err := s.Client.Query(context.Background(), q, variables); err != nil {
			return nil, fmt.Errorf("Failed to fetch additional pull request closing issues references: %w", err)
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

// nolint
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
		"first":  githubv4.Int(githubQuerySize),
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

func (s *Scraper) FetchIssues(owner, name string) ([]*Issue, error) {
	var q struct {
		Repository struct {
			Issues struct {
				Nodes []struct {
					Id          string
					Number      int
					Title       string
					Body        string
					State       string
					StateReason string
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
					// Author
					Author struct {
						Login     string
						AvatarUrl string
					}
				}
				PageInfo PageInfo
			} `graphql:"issues(first: $first, after: $cursor)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	var all []*Issue
	variables := map[string]interface{}{
		"first":  githubv4.Int(githubQuerySize),
		"cursor": (*githubv4.String)(nil),
		"owner":  githubv4.String(owner),
		"name":   githubv4.String(name),
	}

	for {
		if err := s.Client.Query(context.Background(), &q, variables); err != nil {
			return nil, fmt.Errorf("Failed to fetch issues: %w", err)
		}

		if !q.Repository.Issues.PageInfo.HasNextPage {
			break
		}

		issues := make([]*Issue, len(q.Repository.Issues.Nodes))
		for i, is := range q.Repository.Issues.Nodes {
			issue := &Issue{
				Id:          is.Id,
				Number:      is.Number,
				Title:       is.Title,
				Body:        is.Body,
				State:       is.State,
				StateReason: is.StateReason,
				Author:      (*Author)(&is.Author),
				Comments:    nil,
			}

			// Comments
			var cs []*Comment = make([]*Comment, len(is.Comments.Nodes))
			for i, c := range is.Comments.Nodes {
				cs[i] = &Comment{
					Id:     c.Id,
					Body:   c.Body,
					Author: (*Author)(&c.Author),
				}
			}

			if is.Comments.PageInfo.HasNextPage {
				acs, err := s.FetchPullRequestComments(owner, name, is.Number,
					string(is.Comments.PageInfo.EndCursor))
				if err != nil {
					return nil, fmt.Errorf("Failed to fetch issue comments: %w", err)
				}

				cs = append(cs, acs...)
			}

			issue.Comments = cs
			issues[i] = issue
		}
		all = append(all, issues...)

		variables["cursor"] = githubv4.NewString(q.Repository.Issues.PageInfo.EndCursor)
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
					MergedBy       struct {
						Login     string
						AvatarUrl string
					}
					Url string
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
		"first":  githubv4.Int(githubQuerySize),
		"cursor": (*githubv4.String)(nil),
		"owner":  githubv4.String(owner),
		"name":   githubv4.String(name),
	}

	for {
		if err := s.Client.Query(context.Background(), &q, variables); err != nil {
			return nil, fmt.Errorf("Failed to fetch pull requests: %w", err)
		}

		for _, mpr := range q.Repository.PullRequests.Nodes {
			pr := PullRequest{
				Id:             mpr.Id,
				Number:         mpr.Number,
				Title:          mpr.Title,
				Body:           mpr.Body,
				ReviewDecision: mpr.ReviewDecision,
				Merged:         mpr.Merged,
				MergedBy:       (*Author)(&mpr.MergedBy),
				Url:            mpr.Url,
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
					return nil, fmt.Errorf("Failed to fetch pull request closing issue references: %w", err)
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
					return nil, fmt.Errorf("Failed to fetch pull request comments: %w", err)
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
