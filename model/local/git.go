package local

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var (
	ErrCommitEmpty = errors.New("Commit empty")
	ErrBranchEmpty = errors.New("Branch empty")
)

type Hash [20]byte

func (h Hash) ToByte() []byte {
	return h[:]
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
	Content string
}

func NewCommit(o *object.Commit) *Commit {
	if o == nil {
		return nil
	}

	parentHashes := make([]Hash, len(o.ParentHashes))
	for i, hash := range o.ParentHashes {
		parentHashes[i] = Hash(hash)
	}

	return &Commit{
		Hash:         Hash(o.Hash),
		Author:       *NewSignature(&o.Author),
		Committer:    *NewSignature(&o.Committer),
		Message:      o.Message,
		TreeHash:     Hash(o.TreeHash),
		ParentHashes: parentHashes,
	}
}

type Branch struct {
	// Hash of head commit
	Head Hash
	Name string
}

func NewBranch(o *plumbing.Reference) *Branch {
	if o == nil {
		return nil
	}

	return &Branch{
		Head: Hash(o.Hash()),
		Name: o.Name().Short(),
	}
}

type GitModel struct {
	Commits  []Commit
	Branches []Branch
}

func NewGitModel(repo *git.Repository) (*GitModel, error) {
	gitModel := new(GitModel)

	cIter, err := repo.CommitObjects()
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve commits from repository: %w", err)
	}

	err = cIter.ForEach(func(c *object.Commit) error {
		if c == nil {
			return fmt.Errorf("NewGitModel commit: %w", ErrCommitEmpty)
		}
		commit := NewCommit(c)
		gitModel.Commits = append(gitModel.Commits, *commit)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to graft commits to model: %w", err)
	}

	bIter, err := repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve branches from repository: %w", err)
	}
	err = bIter.ForEach(func(b *plumbing.Reference) error {
		if b == nil {
			return fmt.Errorf("NewGitModel branch: %w", ErrBranchEmpty)
		}
		branch := NewBranch(b)
		gitModel.Branches = append(gitModel.Branches, *branch)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to graft branches to model: %w", err)
	}

	return gitModel, nil
}
