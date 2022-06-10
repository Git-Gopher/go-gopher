package detector

import "errors"

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
	Default  string // default branch name
	Branches []struct {
		Name   string       // branch name
		Remote string       // branch remote name
		Head   *CommitGraph // head commit of the branch
	}
}

type FeatureBranchDetect func(branches *FeatureBranchModel) (bool, error)

// FeatureBranchDetector is a detector that detects multiple remote branches (not deleted).
// And check if the branch is a feature branch or a main/develop branch.
type FeatureBranchDetector struct {
	violated int // non feature branches aka develop/release etc. (does not account default branch)
	found    int // total feature branches
	total    int // total branches

	detect FeatureBranchDetect
}

func (bs *FeatureBranchDetector) Run(model *FeatureBranchModel) error {
	if model == nil {
		return ErrFeatureBranchModelNil
	}

	// make sure default branch is always first branch
	// we assume that the default branch likely the main/master branch
	if model.Default != model.Branches[0].Name {
		for i, b := range model.Branches {
			if b.Name == model.Default {
				model.Branches[0], model.Branches[i] = model.Branches[i], model.Branches[0]

				break
			}
		}
	}

	scanBranches := model.Branches

	commitHistory := make(map[*CommitGraph]struct{})

	// Runs BFS on the branches
	for len(scanBranches) != 0 {
		for i, branch := range scanBranches {
			// No more commits
			if branch.Head == nil {
				scanBranches = append(scanBranches[:i], scanBranches[i+1:]...)
				if branch.Name != model.Default {
					// if branch is not default branch, then it is not a feature branch
					bs.violated++
				}

				continue
			}

			// Check if we have already visited this commit
			if _, ok := commitHistory[branch.Head]; ok {
				// Remove scan branch from the list
				scanBranches = append(scanBranches[:i], scanBranches[i+1:]...)

				continue
			}

			// Run if commit has multiple parents and recursively find the common ancestor of the parents
			if len(branch.Head.ParentCommits) > 1 {
				var recursiveCommit func(commits []*CommitGraph) (*CommitGraph, error)
				recursiveCommit = func(commits []*CommitGraph) (*CommitGraph, error) {
					for _, commit := range commits {
						if _, ok := commitHistory[commit]; ok {
							return commit, nil
						}
						commitHistory[commit] = struct{}{}

						if len(commit.ParentCommits) > 1 {
							return recursiveCommit(commit.ParentCommits)
						}
					}
					// No commit that repeats
					return nil, ErrFeatureBranchNoCommonAncestor // should never happen
				}

				nextCommit, err := recursiveCommit(branch.Head.ParentCommits)
				if err != nil {
					return err // should never happen
				}

				branch.Head = nextCommit // set nextCommit as the head of the branch
				scanBranches[i] = branch

				continue
			}

			// Branch has only one parent
			commitHistory[branch.Head] = struct{}{}
			scanBranches[i] = branch
		}
	}

	return nil
}
