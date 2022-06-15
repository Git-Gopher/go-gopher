package github

import (
	"fmt"

	"github.com/Git-Gopher/go-gopher/scrape"
	"github.com/shurcooL/githubv4"
)

type Author struct {
	Login     githubv4.String
	AvatarUrl githubv4.String
}

type Issue struct {
	Id     string
	Body   string
	Title  string
	Author *Author
}

type PullRequest struct {
	Id     string
	Body   string
	Title  string
	Issues []Issue
}

type GithubModel struct {
	Author       *Author
	PullRequests []PullRequest
	Issues       []Issue
}

// TODO: Issues, Author. Also handling the same issue multiple times, should we fetch it multiple
// times or put in memory and search? The former is more memory efficient and is a 'better solution'
// where we can use pointers within our structs, the second is easier in terms of managing complexity
// but also might add complexity in constructing objects multiple times?
func ScrapeGithubModel(owner, name string) (*GithubModel, error) {
	s := scrape.NewScraper()
	sprs, err := s.ScrapePullRequests(owner, name)
	if err != nil {
		return nil, fmt.Errorf("Failed create github model from scraped %w", err)
	}

	// Map scraped PRs to Model PRs
	prs := make([]PullRequest, 0, len(sprs))
	for _, spr := range sprs {
		var issues []Issue
		for _, si := range spr.ClosingIssuesReferences.Edges {
			issues = append(issues, Issue{Id: si.Node.Id, Title: si.Node.Title, Body: si.Node.Body})
		}
		prs = append(prs, PullRequest{Id: spr.Id, Body: spr.Body, Title: spr.Title, Issues: issues})
	}

	return &GithubModel{
		Author:       nil,
		PullRequests: prs,
		Issues:       nil,
	}, nil
}
