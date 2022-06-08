package github

import (
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
//times or put in memory and search? The former is more memory efficient and is a 'better solution'
// where we can use pointers within our structs, the second is easier in terms of managing complexity
// but also might add complexity in constructing objects multiple times?
func ScrapeGithubModel(remote, owner, name string) (*GithubModel, error) {
	s := scrape.NewScraper(remote)
	sprs, err := s.ScrapePullRequests(owner, name)

	if err != nil {
		return nil, err
	}

	// Map Scraped PRs to Model PRs, should this be done on the scrape side or better to keep agnostic?
	var prs []PullRequest
	for _, spr := range sprs.Repository.PullRequests.Nodes {
		var issues []Issue
		for _, sissues := range spr.ClosingIssuesReferences.Edges {
			issues = append(issues, Issue{Id: sissues.Node.Id})
		}
		prs = append(prs, PullRequest{Body: spr.Body, Title: spr.Title, Issues: issues})
	}

	return &GithubModel{
		Author:       nil,
		PullRequests: prs,
		Issues:       nil,
	}, nil
}
