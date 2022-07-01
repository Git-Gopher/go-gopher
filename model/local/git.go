package local

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var (
	ErrCommitEmpty   = errors.New("Commit empty")
	ErrBranchEmpty   = errors.New("Branch empty")
	ErrUnknownLineOp = errors.New("Unknown line op")
)

type Hash [20]byte

func (h Hash) ToByte() []byte {
	return h[:]
}

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

type Signature struct {
	// Name represents a person name. It is an arbitrary string.
	Name string
	// Email is an email, but it cannot be assumed to be well-formed.
	Email string
	// When is the timestamp of the signature.
	When time.Time
}

func NewSignature(o *object.Signature) *Signature {
	if o == nil {
		return nil
	}

	return &Signature{
		Name:  o.Name,
		Email: o.Email,
		When:  o.When,
	}
}

type Diff struct {
	Name     string
	Equal    string
	Addition string
	Deletion string

	Points []DiffPoint
}

type DiffPoint struct {
	OldPosition int64
	OldLines    int64

	NewPosition int64
	NewLines    int64

	LinesAdded   int64
	LinesDeleted int64
}

type File = gitdiff.File

type Commit struct {
	// Hash of the commit object.
	Hash Hash
	// TreeHash is the hash of the root tree of the commit.
	TreeHash Hash
	// ParentHashes are the hashes of the parent commits of the commit.
	ParentHashes []Hash
	// Author is the original author of the commit.
	Author Signature
	// Committer is the one performing the commit, might be different from Author.
	Committer Signature
	// Message is the commit message, contains arbitrary text.
	Message string
	// TODO: Import go-git types
	Content       string
	DiffToParents []Diff
}

type Operation int

const (
	// Equal item represents a equals diff.
	Equal Operation = iota
	// Add item represents an insert diff.
	Add
	// Delete item represents a delete diff.
	Delete
)

type Chunk struct {
	// Content contains the portion of the file.
	Content string
	// Type contains the Operation to do with this Chunk.
	Type Operation
}

func NewCommit(r *git.Repository, c *object.Commit) *Commit {
	if c == nil || r == nil {
		return nil
	}

	parentHashes := make([]Hash, len(c.ParentHashes))
	for i, hash := range c.ParentHashes {
		parentHashes[i] = Hash(hash)
	}

	var diffs []Diff
	err := c.Parents().ForEach(
		func(p *object.Commit) error {
			diff, err := FetchDiffs(p, c)
			if err != nil {
				return fmt.Errorf("Failed to fetch diff: %w", err)
			}
			diffs = append(diffs, diff...)

			return nil
		})
	if err != nil {
		return nil
	}

	return &Commit{
		Hash:          Hash(c.Hash),
		Author:        *NewSignature(&c.Author),
		Committer:     *NewSignature(&c.Committer),
		Message:       c.Message,
		TreeHash:      Hash(c.TreeHash),
		ParentHashes:  parentHashes,
		DiffToParents: diffs,
	}
}

// TODO: Might be useful to add some of these to the Branch struct.
// type MockBranchModel struct {
// 	Ref           string
// 	Remote        string
// 	Hash          string
// 	CommitsBehind int       // Number of commits behind the primary branch
// 	LastChange    time.Time // Time of the head commit of the current branch
// }.
type Branch struct {
	// Hash of head commit
	Head Commit
	Name string
}

func NewBranch(repo *git.Repository, o *plumbing.Reference, c *object.Commit) *Branch {
	if o == nil {
		return nil
	}

	return &Branch{
		Head: *NewCommit(repo, c),
		Name: o.Name().Short(),
	}
}

type GitModel struct {
	Commits      []Commit
	Branches     []Branch
	MainGraph    *BranchGraph
	BranchMatrix []*BranchMatrix
}

func NewGitModel(repo *git.Repository) (*GitModel, error) {
	gitModel := new(GitModel)

	// Commits
	cIter, err := repo.CommitObjects()
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve commits from repository: %w", err)
	}
	err = cIter.ForEach(func(c *object.Commit) error {
		if c == nil {
			return fmt.Errorf("NewGitModel commit: %w", ErrCommitEmpty)
		}
		commit := NewCommit(repo, c)
		gitModel.Commits = append(gitModel.Commits, *commit)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to graft commits to model: %w", err)
	}

	// Branches
	bIter, err := repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve branches from repository: %w", err)
	}
	err = bIter.ForEach(func(b *plumbing.Reference) error {
		if b == nil {
			return fmt.Errorf("NewGitModel branch: %w", ErrBranchEmpty)
		}

		var c *object.Commit
		c, err = repo.CommitObject(b.Hash())
		if err != nil {
			return fmt.Errorf("Failed to find head commit from branch: %w", err)
		}

		branch := NewBranch(repo, b, c)
		gitModel.Branches = append(gitModel.Branches, *branch)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to graft branches to model: %w", err)
	}

	// MainGraph
	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("Failed to find head reference: %w", err)
	}
	refCommit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("Failed to find head commit: %w", err)
	}
	gitModel.MainGraph = FetchBranchGraph(refCommit)
	gitModel.MainGraph.BranchName = ref.Hash().String()

	// BranchMatrix
	branches := []plumbing.Hash{}
	bIter, err = repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve branches from repository: %w", err)
	}
	_ = bIter.ForEach(func(b *plumbing.Reference) error {
		if b == nil {
			return nil
		}
		branches = append(branches, b.Hash())

		return nil
	})

	gitModel.BranchMatrix, err = CreateBranchMatrix(repo, branches)
	if err != nil {
		return nil, fmt.Errorf("Failed to create branch matrix: %w", err)
	}

	return gitModel, nil
}
