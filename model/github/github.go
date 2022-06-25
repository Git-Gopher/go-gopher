package github

import "fmt"

type Author struct {
	Login     string
	AvatarUrl string
}

type Issue struct {
	Id          string
	Number      int
	Title       string
	Body        string
	State       string
	StateReason string
	Author      *Author
	Comments    []*Comment
}

type Comment struct {
	Id     string
	Body   string
	Author *Author
}

type ReviewThread struct {
	Id         string
	IsResolved bool
	IsOutdated bool
	Path       string
}

type PullRequest struct {
	Id             string
	Number         int
	Title          string
	Body           string
	ReviewDecision string
	Merged         bool
	MergedBy       *Author
	Url            string
	Author         *Author
	ClosingIssues  []*Issue
	Comments       []*Comment
	ReviewThreads  []*ReviewThread
}

type GithubModel struct {
	Author       *Author
	PullRequests []*PullRequest
	Issues       []*Issue
}

// TODO: Issues, Author. Also handling the same issue multiple times, should we fetch it multiple
// times or put in memory and search? The former is more memory efficient and is a 'better solution'
// where we can use pointers within our structs, the second is easier in terms of managing complexity
// but also might add complexity in constructing objects multiple times?
func ScrapeGithubModel(owner, name string) (*GithubModel, error) {
	s := NewScraper()
	prs, err := s.FetchPullRequests(owner, name)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch pull requests for GitHub model: %w", err)
	}

	issues, err := s.FetchIssues(owner, name)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch issues for GitHub model: %w", err)
	}

	thing := GithubModel{
		Author:       nil,
		PullRequests: prs,
		Issues:       issues,
	}

	return &thing, nil
}
