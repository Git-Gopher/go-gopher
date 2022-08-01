package analysis

import (
	"github.com/Git-Gopher/go-gopher/model/enriched"
)

// contribution.go contains code to find the contribution of a user to a detector.

// contribution contains maps of contribution counts for a user by email.
type contribution struct {
	BranchCountMap      map[string]int
	CommitCountMap      map[string]int
	PullRequestCountMap map[string]int
	MergeCountMap       map[string]int
}

// NewContribution calculates the contribution by email.
func NewContribution(enriched enriched.EnrichedModel) *contribution {
	c := &contribution{
		BranchCountMap:      make(map[string]int),
		CommitCountMap:      make(map[string]int),
		PullRequestCountMap: make(map[string]int),
		MergeCountMap:       make(map[string]int),
	}

	for _, branch := range enriched.Branches {
		email := branch.Head.Author.Email
		if _, ok := c.BranchCountMap[email]; !ok {
			c.BranchCountMap[email] = 0
		}
		c.BranchCountMap[email]++

		if email == branch.Head.Committer.Email {
			continue
		}

		email = branch.Head.Committer.Email
		if _, ok := c.BranchCountMap[email]; !ok {
			c.BranchCountMap[email] = 0
		}
		c.BranchCountMap[email]++
	}

	for _, commit := range enriched.Commits {
		email := commit.Author.Email
		if _, ok := c.CommitCountMap[email]; !ok {
			c.CommitCountMap[email] = 0
		}
		c.CommitCountMap[email]++

		if email == commit.Committer.Email {
			continue
		}

		email = commit.Committer.Email
		if _, ok := c.CommitCountMap[email]; !ok {
			c.CommitCountMap[email] = 0
		}
		c.CommitCountMap[email]++
	}

	for _, pr := range enriched.PullRequests {
		email := pr.Author.Email
		if _, ok := c.PullRequestCountMap[email]; !ok {
			c.PullRequestCountMap[email] = 0
		}
		c.PullRequestCountMap[email]++
	}

	for _, commit := range enriched.Commits {
		if len(commit.ParentHashes) < 2 {
			continue
		}

		email := commit.Author.Email
		if _, ok := c.MergeCountMap[email]; !ok {
			c.MergeCountMap[email] = 0
		}
		c.MergeCountMap[email]++

		if email == commit.Committer.Email {
			continue
		}

		email = commit.Committer.Email
		if _, ok := c.MergeCountMap[email]; !ok {
			c.MergeCountMap[email] = 0
		}
		c.MergeCountMap[email]++
	}

	return c
}
