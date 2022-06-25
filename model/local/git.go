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
	// Patch
	PatchChucks []map[string][]Chunk
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

func FetchChunk(from *object.Commit, to *object.Commit) ([]Chunk, error) {
	patch, err := from.Patch(to)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch chunk: %w", err)
	}

	chunks := []Chunk{}

	fmt.Println("++++++++++++++++++")
	fmt.Println(to.Hash.String())
	fmt.Println(from.Hash.String())
	fmt.Println("++++++++++++++++++")
	fmt.Println(patch.String())

	for _, filePatch := range patch.FilePatches() {
		fileChunks := filePatch.Chunks()
		fromFile, toFile := filePatch.Files()
		fmt.Printf("from file: %v\n", fromFile.Path())
		fmt.Printf("to file: %v\n", toFile.Path())
		for _, chunk := range fileChunks {
			chunk.Content()
		}
		// if err != nil {
		// 	return nil, fmt.Errorf("failed to fetch filepatch: %w", err)
		// }
		// fmt.Printf("files: %v\n", files)
		for _, chunk := range fileChunks {

			fmt.Printf("chunk: %v\n", chunk.Content())
			chunks = append(chunks, Chunk{
				Content: chunk.Content(),
				Type:    Operation(chunk.Type()),
			})
		}
	}

	return chunks, nil
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
