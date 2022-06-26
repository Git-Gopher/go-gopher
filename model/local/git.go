package local

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
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

type Diff struct {
	Name     string
	Equal    string
	Addition string
	Deletion string
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

func FetchDiffs(from *object.Commit, to *object.Commit) ([]Diff, error) {
	patch, err := from.Patch(to)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch chunk: %w", err)
	}

	files, _, err := gitdiff.Parse(strings.NewReader(patch.String()))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse diff: %w", err)
	}
	diffs := make([]Diff, len(files))
	for i, f := range files {
		equal, added, deleted, err := Defragment(f.TextFragments)
		if err != nil {
			return nil, fmt.Errorf("Failed to defragment diffs: %w", err)
		}

		var name string
		if f.IsNew || f.IsRename {
			name = f.NewName
		} else {
			name = f.OldName
		}

		diffs[i] = Diff{
			Name:     name,
			Addition: added,
			Deletion: deleted,
			Equal:    equal,
		}
	}

	return diffs, nil
}

func Defragment(fragment []*gitdiff.TextFragment) (equal, added, deleted string, err error) {
	for _, f := range fragment {
		for _, l := range f.Lines {
			switch l.Op {
			case gitdiff.LineOp(Equal):
				equal += l.Line
			case gitdiff.LineOp(Add):
				added += l.Line
			case gitdiff.LineOp(Delete):
				deleted += l.Line
			default:
				return "", "", "", fmt.Errorf("Unexpected Op: %v", l.Op)
			}
		}
	}

	return equal, added, deleted, nil
}

func NewCommit(c *object.Commit) *Commit {
	if c == nil {
		return nil
	}

	parentHashes := make([]Hash, len(c.ParentHashes))
	for i, hash := range c.ParentHashes {
		parentHashes[i] = Hash(hash)
	}

	return &Commit{
		Hash:         Hash(c.Hash),
		Author:       *NewSignature(&c.Author),
		Committer:    *NewSignature(&c.Committer),
		Message:      c.Message,
		TreeHash:     Hash(c.TreeHash),
		ParentHashes: parentHashes,
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

func NewBranch(o *plumbing.Reference, c *object.Commit) *Branch {
	if o == nil {
		return nil
	}

	return &Branch{
		Head: *NewCommit(c),
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

		var c *object.Commit
		c, err = repo.CommitObject(b.Hash())
		if err != nil {
			return fmt.Errorf("Failed to find head commit from branch: %w", err)
		}

		branch := NewBranch(b, c)
		gitModel.Branches = append(gitModel.Branches, *branch)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to graft branches to model: %w", err)
	}

	return gitModel, nil
}
