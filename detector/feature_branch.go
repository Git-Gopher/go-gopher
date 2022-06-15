package detector

import (
	"errors"
	"math"

	"github.com/Git-Gopher/go-gopher/violation"
)

var (
	ErrFeatureBranchModelNil         = errors.New("feature branch model is nil")
	ErrFeatureBranchNoCommonAncestor = errors.New("feature branch no common ancestor found")
)

// TODO: move to models
type CommitGraph struct {
	Hash          string
	ParentCommits []*CommitGraph
}

type FeatureBranchModel struct {
	BranchName string
	Head       *CommitGraph
}

type FeatureBranchDetect func(branches *FeatureBranchModel) (bool, error)

// FeatureBranchDetector is a detector that detects multiple remote branches (not deleted).
// And check if the branch is a feature branch or a main/develop branch.
type FeatureBranchDetector struct {
	violated   int // non feature branches aka develop/release etc. (does not account default branch)
	found      int // total feature branches
	total      int // total branches
	violations []violation.Violation

	primaryBranch string

	detect FeatureBranchDetect
}

func (bs *FeatureBranchDetector) Run(model *FeatureBranchModel) error {
	if model == nil {
		return ErrFeatureBranchModelNil
	}

	bs.violated = 0
	bs.found = 0
	bs.total = 0
	bs.violations = make([]violation.Violation, 0)

	bs.primaryBranch = model.BranchName

	return nil
}

func (bs *FeatureBranchDetector) checkNext(c *CommitGraph) *CommitGraph {
	if c == nil {
		return nil
	}

	// if it has multiple parents
	if len(c.ParentCommits) > 1 {
		for _, child := range c.ParentCommits {
			if len(child.ParentCommits) > 1 {
				// skip other branch
				return bs.checkNext(child)
			}
		}

		// if both parents have one commit
		// assumes the shorter branch to be the violation
		lenViolations := math.MaxInt
		var violations []violation.Violation
		var nextCommit *CommitGraph

		for _, child := range c.ParentCommits {
			next, v := bs.checkEnd(child, []violation.Violation{})
			if next == nil {
				// no more commits to check
				bs.violations = append(bs.violations, v...)

				return nil
			}
			if len(v) < lenViolations {
				lenViolations = len(v)
				violations = v
				nextCommit = next
			}
		}

		bs.violations = append(bs.violations, violations...)

		return bs.checkNext(nextCommit)
	}

	// only one parent (violation)
	bs.violations = append(bs.violations, violation.NewPrimaryBranchDirectCommitViolation(
		bs.primaryBranch,
		c.Hash,
		[]string{c.ParentCommits[0].Hash},
	))

	return bs.checkNext(c.ParentCommits[0])

}

// Check if the commit is made as the start of the branch
// if not return last commit with two parent and associated violations
func (bs *FeatureBranchDetector) checkEnd(
	c *CommitGraph,
	v []violation.Violation,
) (*CommitGraph, []violation.Violation) {
	if c == nil {
		return nil, v
	}

	// No more commits to recursive check
	if len(c.ParentCommits) == 0 {
		// All violations are removed as the start of the branch
		return nil, []violation.Violation{}
	}

	// The parent has two commits
	if len(c.ParentCommits) > 1 {
		return c, v
	}

	// The parent has one commit
	v = append(v, violation.NewPrimaryBranchDirectCommitViolation(
		bs.primaryBranch,
		c.Hash,
		[]string{c.ParentCommits[0].Hash},
	))

	return bs.checkEnd(c, v)
}
