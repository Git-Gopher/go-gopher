package detector

import (
	"errors"
	"math"

	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/violation"
)

var (
	ErrFeatureBranchModelNil         = errors.New("feature branch model is nil")
	ErrFeatureBranchNoCommonAncestor = errors.New("feature branch no common ancestor found")
)

// NewFeatureBranchDetector creates a new feature branch detector.
// This detector is a custom recursive detector that does not rely on the
// `detector.NewDetector(detector.Detect)` pattern.
func NewFeatureBranchDetector() *FeatureBranchDetector {
	return &FeatureBranchDetector{
		violated:   0,
		found:      0,
		total:      0,
		violations: make([]violation.Violation, 0),
	}
}

// FeatureBranchDetector is a detector that detects multiple remote branches (not deleted).
// And check if the branch is a feature branch or a main/develop branch.
type FeatureBranchDetector struct {
	violated   int // non feature branches aka develop/release etc. (does not account default branch)
	found      int // total feature branches
	total      int // total branches
	violations []violation.Violation

	primaryBranch string
}

func (bs *FeatureBranchDetector) Run(model *enriched.EnrichedModel) error {
	if model == nil {
		return ErrFeatureBranchModelNil
	}

	bs.violated = 0
	bs.found = 0
	bs.total = 0
	bs.violations = make([]violation.Violation, 0)

	bs.primaryBranch = model.MainGraph.BranchName

	bs.checkNext(model.MainGraph.Head)

	return nil
}

// checkNext is used to check the next commit in the branch. (recursive)
//
// If the commit has multiple parents, it will check which is the primary parent.
// If the commit has one parent check if all commits after this commit has one parent.
// If the commit has one parent and one parent after it has multiple parents,
// this commit is a direct commit not from a feature branch (violation).
func (bs *FeatureBranchDetector) checkNext(c *local.CommitGraph) *local.CommitGraph {
	if c == nil {
		return nil
	}

	bs.total += 1

	// if it has no parents
	if len(c.ParentCommits) == 0 {
		return nil
	}

	// if it has multiple parents
	if len(c.ParentCommits) > 1 {
		for _, child := range c.ParentCommits {
			if len(child.ParentCommits) > 1 {
				bs.found += 1
				// skip other branch
				return bs.checkNext(child)
			}
		}

		// if both parents have one commit
		// assumes the shorter branch to be the violation
		lenViolations := math.MaxInt
		var violations []violation.Violation
		var nextCommit *local.CommitGraph

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

	next, v := bs.checkEnd(c.ParentCommits[0], []violation.Violation{})
	if next == nil {
		// no more commits to check
		bs.violations = append(bs.violations, v...)

		return nil
	}

	// only one parent (violation)
	bs.violations = append(bs.violations, violation.NewPrimaryBranchDirectCommitViolation(
		bs.primaryBranch,
		c.Hash,
		[]string{c.ParentCommits[0].Hash},
	))
	bs.violated++

	return bs.checkNext(c.ParentCommits[0])
}

// Check if the commit is made as the start of the branch
// if not return last commit with two parent and associated violations.
func (bs *FeatureBranchDetector) checkEnd(
	c *local.CommitGraph,
	v []violation.Violation,
) (*local.CommitGraph, []violation.Violation) {
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
	bs.violated++

	return bs.checkEnd(c, v)
}

func (bs *FeatureBranchDetector) Result() (int, int, int, []violation.Violation) {
	return bs.violated, bs.found, bs.total, bs.violations
}
