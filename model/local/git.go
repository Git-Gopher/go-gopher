package local

import (
	"crypto/md5" //nolint:gosec // we don't need a strong hash here
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	log "github.com/sirupsen/logrus"
)

var (
	ErrCommitEmpty     = errors.New("commit empty")
	ErrBranchEmpty     = errors.New("branch empty")
	ErrUnknownLineOp   = errors.New("unknown line op")
	ErrBadTagReference = errors.New("empty reference, repo or commit for tag")
	ErrNewCommitNil    = errors.New("commit is nil")
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
	// PatchID is the hash of the patch. If empty means more than one parent (not cherry-picked)
	PatchID *string `json:"-"`
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

func NewCommit(r *git.Repository, c *object.Commit) (*Commit, error) {
	if c == nil || r == nil {
		return nil, fmt.Errorf("%w: %v, %v", ErrNewCommitNil, c, r)
	}

	parentHashes := []Hash{}
	// Root commit has empty commit tree as parent. This also include orphan commits.
	if len(c.ParentHashes) == 0 {
		h, err := hex.DecodeString(EmptyTreeHash)
		if err != nil {
			return nil, fmt.Errorf("unable to decode emptytreehash to a byte hash: %s, %w", EmptyTreeHash, err)
		}

		var arr [20]byte
		copy(arr[:], h[:20])
		parentHashes = append(parentHashes, Hash(arr))
	} else {
		for _, hash := range c.ParentHashes {
			parentHashes = append(parentHashes, Hash(hash))
		}
	}

	var diffs []Diff
	commitTree, err := c.Tree()
	if err != nil {
		return nil, fmt.Errorf("cannot fetch commit tree: %w", err)
	}

	parentTree := &object.Tree{}
	var patch *object.Patch
	if c.NumParents() != 0 { // nolint: nestif
		var parent *object.Commit
		parent, err = c.Parents().Next()
		if err != nil {
			return nil, fmt.Errorf("cannot fetch commit parents: %w", err)
		}
		parentTree, err = parent.Tree()
		if err != nil {
			return nil, fmt.Errorf("cannot fetch commit parent tree: %w", err)
		}
		patch, err = parentTree.Patch(commitTree)
		if err != nil {
			return nil, fmt.Errorf("cannot fetch commit patch: %w", err)
		}
	} else {
		patch, err = parentTree.Patch(commitTree)
		if err != nil {
			return nil, fmt.Errorf("cannot fetch commit patch: %w", err)
		}
	}

	var patchID *string

	// calculate patch id
	h := md5.New() //nolint:gosec // we don't need a strong hash here

	patchDiff := strings.Builder{}
	for _, diff := range diffs {
		patchDiff.WriteString(diff.Name)
		patchDiff.WriteString(diff.Equal)
		patchDiff.WriteString(diff.Addition)
		patchDiff.WriteString(diff.Deletion)
	}
	io.WriteString(h, patchDiff.String()) //nolint:errcheck,gosec
	p := fmt.Sprintf("%x", h.Sum(nil))
	patchID = &p

	// diffs
	diff, err := FetchDiffs(patch)
	if err != nil {
		return nil, fmt.Errorf("cannot commit diffs: %w", err)
	}

	diffs = append(diffs, diff...)

	return &Commit{
		Hash:          Hash(c.Hash),
		Author:        *NewSignature(&c.Author),
		Committer:     *NewSignature(&c.Committer),
		Message:       c.Message,
		TreeHash:      Hash(c.TreeHash),
		ParentHashes:  parentHashes,
		DiffToParents: diffs,
		PatchID:       patchID,
	}, nil
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
	commit, err := NewCommit(repo, c)
	if err != nil {
		log.Warnf("unable to find commit for branch: %v", err)
	}

	return &Branch{
		Head: *commit,
		Name: strings.Replace(o.Name().Short(), "origin/", "", 1),
	}
}

type Tag struct {
	// Name of the tag. Eg: v0.0.8.
	Name string
	// Head of the tag.
	Head Commit
}

func NewTag(repo *git.Repository, o *plumbing.Reference, c *object.Commit) (*Tag, error) {
	if o == nil || c == nil || repo == nil {
		return nil, fmt.Errorf("%w: %v, %v, %v", ErrBadTagReference, repo, o, c)
	}

	commit, err := NewCommit(repo, c)
	if err != nil {
		return nil, fmt.Errorf("unable to find commit for tag: %w", err)
	}

	return &Tag{
		Name: o.Name().Short(),
		Head: *commit,
	}, nil
}

type GitModel struct {
	Commits      []Commit
	Committer    []Committer
	Branches     []Branch
	MainGraph    *BranchGraph
	BranchMatrix []*BranchMatrix
	Tags         []*Tag

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

		var commit *Commit
		commit, err = NewCommit(repo, c)
		if err != nil {
			return fmt.Errorf("unable to find commit while creating git model: %w", err)
		}

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

	// Branches
	branches := []plumbing.Hash{}
	rIter, err := repo.References()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve branches from repository: %w", err)
	}
	err = rIter.ForEach(func(b *plumbing.Reference) error {
		if !b.Name().IsRemote() {
			// not a branch
			return nil
		}

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

		if b.Hash().String() == ref.Hash().String() {
			gitModel.MainGraph.BranchName = strings.Replace(ref.Name().Short(), "origin/", "", 1)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to graft branches to model: %w", err)
	}

	// BranchMatrix
	gitModel.BranchMatrix, err = CreateBranchMatrix(repo, branches)
	if err != nil {
		return nil, fmt.Errorf("failed to create branch matrix: %w", err)
	}

	// Tags.
	tagIter, err := repo.Tags()
	if err != nil {
		return nil, fmt.Errorf("failed to generate tag iter: %w", err)
	}

	var ts []*Tag
	if err = tagIter.ForEach(func(o *plumbing.Reference) error {
		if o == nil {
			return fmt.Errorf("nil tag reference: %w", err)
		}

		var c *object.Commit
		c, err = repo.CommitObject(o.Hash())
		if err != nil {
			return fmt.Errorf("failed to find head commit from branch: %w", err)
		}

		var t *Tag
		t, err = NewTag(repo, o, c)
		if err != nil {
			log.Warnf("Unable to create tag: %v", err)

			return err
		}
		ts = append(ts, t)

		return nil
	}); err != nil {
		return nil, fmt.Errorf("bad tag iteration: %w", err)
	}

	gitModel.Tags = ts

	gitModel.Repository = repo

	return gitModel, nil
}
