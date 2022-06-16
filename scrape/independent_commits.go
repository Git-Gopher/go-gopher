package scrape

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func CreateIndependentCommits(
	r *git.Repository,
	branchAHead *object.Commit,
	branchBHead *object.Commit,
) ([]*object.Commit, error) {
	independent, err := object.Independents([]*object.Commit{branchAHead, branchBHead})
	if err != nil {
		return nil, fmt.Errorf("failed to create independent commits: %w", err)
	}

	return independent, nil
}
