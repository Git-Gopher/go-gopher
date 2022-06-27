package local

import (
	"github.com/go-git/go-git/v5/plumbing/object"
)

// BranchGraph represents a branch as a CommitGraph
type BranchGraph struct {
	// Branch Name of the represented branch
	BranchName string
	// The first commit represented as a CommitGraph
	Head *CommitGraph
}

// CommitGraph represents a commit graph based on a hash and its parent commits.
type CommitGraph struct {
	// Hash of the commit
	Hash string
	// Parent commits represented as a CommitGraph
	ParentCommits []*CommitGraph
}

func FetchBranchGraph(head *object.Commit) *BranchGraph {
	// new BranchGraph
	branchGraph := new(BranchGraph)

	// stack of commits to process
	// commitRefs[i] and commitGraphRef[i] are correlated
	commitRefs := make([]*object.Commit, 0)
	commitGraphRefs := make([]*CommitGraph, 0)

	// caches previous commits
	commitGraphMap := make(map[string]*CommitGraph)

	// head of the graph
	hash := head.Hash.String()
	headGraph := &CommitGraph{Hash: hash}
	branchGraph.Head = headGraph

	err := head.Parents().ForEach(
		func(c *object.Commit) error {
			hash := c.Hash.String()
			commit := &CommitGraph{Hash: hash}

			// add commit to stack and cache
			commitRefs = append(commitRefs, c)
			commitGraphRefs = append(commitGraphRefs, commit)
			commitGraphMap[hash] = commit

			// add commit to graph
			headGraph.ParentCommits = append(headGraph.ParentCommits, commit)

			return nil
		})
	if err != nil {
		return nil
	}

	// process commits in the stack
	for len(commitRefs) != 0 {
		// pop from the stack
		// pop last commit for efficiency
		n := len(commitRefs) - 1
		commit := commitRefs[n]
		commitGraph := commitGraphRefs[n]

		// remove last commit from list
		commitRefs = commitRefs[:n]
		commitGraphRefs = commitGraphRefs[:n]

		_ = commit.Parents().ForEach(
			func(c *object.Commit) error {
				hash := commit.Hash.String()

				// check cache if exist use commit object and check next
				if cached, ok := commitGraphMap[hash]; !ok {
					commitGraph.ParentCommits = append(commitGraph.ParentCommits, cached)

					return nil
				}

				// add parent to stack to be processed
				commit := &CommitGraph{Hash: hash}
				commitRefs = append(commitRefs, c)
				commitGraphRefs = append(commitGraphRefs, commit)
				commitGraphMap[hash] = commit

				return nil
			})
	}

	return branchGraph
}
