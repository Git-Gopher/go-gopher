package remote

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Author struct {
	Login     string
	AvatarUrl string
	Email     string
}

type Issue struct {
	Id          string
	Number      int
	Title       string
	Body        string
	State       string
	StateReason string
	Author      *Author
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
	HeadRefName    string // source branch
	BaseRefName    string // target branch
	CreatedAt      *time.Time
	ClosedAt       *time.Time
	Title          string
	Body           string
	ReviewDecision string
	Merged         bool
	MergedBy       *Author
	Url            string
	Author         *Author
	ClosingIssues  []*Issue
	ReviewThreads  []*ReviewThread
}

type RemoteModel struct {
	Owner        string
	Name         string
	URL          string
	PullRequests []*PullRequest
	Issues       []*Issue
	Committers   []Committer
}

type Committer struct {
	CommitId string
	Email    string
	Login    string
}

// TODO: Issues, Author. Also handling the same issue multiple times, should we fetch it multiple
// times or put in memory and search? The former is more memory efficient and is a 'better solution'
// where we can use pointers within our structs, the second is easier in terms of managing complexity
// but also might add complexity in constructing objects multiple times?
func ScrapeRemoteModel(owner, name string) (*RemoteModel, error) {
	ghm := RemoteModel{
		Owner:        owner,
		Name:         name,
		URL:          "",
		PullRequests: nil,
		Issues:       nil,
		Committers:   nil,
	}

	s := NewScraper()

	var err error
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error)
	waitCh := make(chan struct{})

	wg.Add(1)
	go func() {
		ghm.PullRequests, err = s.FetchPullRequests(ctx, owner, name)
		if err != nil {
			errCh <- fmt.Errorf("Failed to fetch issues for GitHub model: %w", err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		ghm.Issues, err = s.FetchIssues(ctx, owner, name)
		if err != nil {
			errCh <- fmt.Errorf("Failed to fetch issues for GitHub model: %w", err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		ghm.URL, err = s.FetchURL(ctx, owner, name)
		if err != nil {
			errCh <- fmt.Errorf("Failed to fetch issues for GitHub model: %w", err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		ghm.Committers, err = s.FetchCommitters(ctx, owner, name)
		if err != nil {
			errCh <- fmt.Errorf("Failed to fetch issues for GitHub model: %w", err)
		}
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(waitCh)
	}()

	select {
	case <-waitCh:
		// All done
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
	}

	return &ghm, nil
}
