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
	ErrCommitEmpty   = errors.New("commit empty")
	ErrBranchEmpty   = errors.New("branch empty")
	ErrUnknownLineOp = errors.New("unknown line op")
	// Hash of an empty git tree.
	// $(printf '' | git hash-object -t tree --stdin).
	EmptyTreeHash = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"
)

type Hash [20]byte

func (h Hash) ToByte() []byte {
	return h[:]
}

func (h Hash) String() string {
	return string(h[:])
}

func (h Hash) HexString() string {
	return hex.EncodeToString((h[:]))
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
	IsBinary bool
	Equal    string
	Addition string
	Deletion string

	Points []DiffPoint `json:"-"`
}

type DiffPoint struct {
	OldPosition int64
	NewPosition int64

	LinesAdded   int64
	LinesDeleted int64
}

type File = gitdiff.File

type Commit struct {
	// Hash of the commit object.
	Hash Hash
	// TreeHash is the hash of the root tree of the commit.
	TreeHash Hash `json:"-"`
	// ParentHashes are the hashes of the parent commits of the commit.
	ParentHashes []Hash `json:"-"`
	// Author is the original author of the commit.
	Author Signature
	// Committer is the one performing the commit, might be different from Author.
	Committer Signature `json:"-"`
	// Message is the commit message, contains arbitrary text.
	Message       string
	Content       string
	DiffToParents []Diff `json:"-"`
}

type Committer struct {
	CommitId string
	Email    string
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
	if len(parentHashes) == 0 { //nolint: nestif
		iter, err := r.TreeObjects()
		if err != nil {
			return nil
		}
		if err := iter.ForEach(func(o *object.Tree) error {
			changes, err := o.Diff(&object.Tree{Hash: plumbing.NewHash(EmptyTreeHash)})
			if err != nil {
				return fmt.Errorf("failed to fetch tree root diff: %w", err)
			}

			patch, err := changes.Patch()
			if err != nil {
				return fmt.Errorf("failed to fetch root patch: %w", err)
			}

			diff, err := FetchDiffs(patch)
			if err != nil {
				return fmt.Errorf("failed to fetch root diff: %w", err)
			}
			diffs = append(diffs, diff...)

			return err
		}); err != nil {
			return nil
		}
	} else {
		err := c.Parents().ForEach(
			func(p *object.Commit) error {
				var diff []Diff
				var patch *object.Patch
				patch, err := p.Patch(c)
				if err != nil {
					return fmt.Errorf("failed to fetch patch: %w", err)
				}

				diff, err = FetchDiffs(patch)
				if err != nil {
					return fmt.Errorf("failed to fetch diff: %w", err)
				}
				diffs = append(diffs, diff...)

				return nil
			})
		if err != nil {
			return nil
		}
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
//
//	type MockBranchModel struct {
//		Ref           string
//		Remote        string
//		Hash          string
//		CommitsBehind int       // Number of commits behind the primary branch
//		LastChange    time.Time // Time of the head commit of the current branch
//	}.
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
	Committer    []Committer
	Branches     []Branch
	MainGraph    *BranchGraph
	BranchMatrix []*BranchMatrix

	// Not all functionality has been ported from go-git.
	Repository *git.Repository
}

func NewGitModel(repo *git.Repository) (*GitModel, error) {
	gitModel := new(GitModel)

	// Commits
	cIter, err := repo.CommitObjects()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve commits from repository: %w", err)
	}
	err = cIter.ForEach(func(c *object.Commit) error {
		if c == nil {
			return fmt.Errorf("NewGitModel commit: %w", ErrCommitEmpty)
		}
		commit := NewCommit(repo, c)
		gitModel.Commits = append(gitModel.Commits, *commit)
		gitModel.Committer = append(gitModel.Committer, Committer{
			CommitId: string(c.Hash[:]),
			Email:    c.Committer.Email,
		})

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to graft commits to model: %w", err)
	}

	// Branches
	branches := []plumbing.Hash{}
	bIter, err := repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve branches from repository: %w", err)
	}
	err = bIter.ForEach(func(b *plumbing.Reference) error {
		if b == nil {
			return fmt.Errorf("NewGitModel branch: %w", ErrBranchEmpty)
		}
		branches = append(branches, b.Hash())

		var c *object.Commit
		c, err = repo.CommitObject(b.Hash())
		if err != nil {
			return fmt.Errorf("failed to find head commit from branch: %w", err)
		}

		branch := NewBranch(repo, b, c)
		gitModel.Branches = append(gitModel.Branches, *branch)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to graft branches to model: %w", err)
	}

	// MainGraph
	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to find head reference: %w", err)
	}
	refCommit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to find head commit: %w", err)
	}
	gitModel.MainGraph = FetchBranchGraph(refCommit)
	gitModel.MainGraph.BranchName = ref.Hash().String()

	// BranchMatrix
	gitModel.BranchMatrix, err = CreateBranchMatrix(repo, branches)
	if err != nil {
		return nil, fmt.Errorf("failed to create branch matrix: %w", err)
	}

	gitModel.Repository = repo

	return gitModel, nil
}
