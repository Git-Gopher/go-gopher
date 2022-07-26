package detector

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/model/remote"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

func TestTwoParentsCommitDetect(t *testing.T) {
	r := utils.FetchRepository(t, "https://github.com/Git-Gopher/tests", "test/two-parents-merged/0")
	tests := []struct {
		name string
		want CommitDetect
	}{
		{"asdf", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the gitModel
			gitModel, err := local.NewGitModel(r)
			if err != nil {
				t.Errorf("TestTwoParentsCommitDetect() create model = %v", err)
			}

			enrichedModel := enriched.NewEnrichedModel(*gitModel, remote.RemoteModel{})

			detector := NewCommitDetector(BranchCommitDetect())
			if err = detector.Run(enrichedModel); err != nil {
				t.Errorf("TestTwoParentsCommitDetect() run detector = %v", err)
			}

			t.Log(detector.Result())
		})
	}
}

func TestTwoParentsCommitDetectGoGit(t *testing.T) {
	// Setup go git repo with the configuration that we want (two parents one commit)
	fs := memfs.New()
	storer := memory.NewStorage()

	r, err := git.Init(storer, fs)
	if err != nil {
		t.Errorf("TestTwoParentsCommitDetectGoGit() init repo = %v", err)
	}

	w, err := r.Worktree()
	if err != nil {
		t.Errorf("TestTwoParentsCommitDetectGoGit() worktree = %v", err)
	}

	file, err := fs.Create("main.go")
	if err != nil {
		t.Errorf("TestTwoParentsCommitDetectGoGit() create file main.go = %v", err)
	}

	_, err = file.Write([]byte(`package main\n\nimport "fmt"\n\nfunc main() {\n	fmt.Println("Hello World")\n}\n`))
	if err != nil {
		t.Errorf("TestTwoParentsCommitDetectGoGit() write = %v", err)
	}

	err = file.Close()
	if err != nil {
		t.Errorf("TestTwoParentsCommitDetectGoGit() close = %v", err)
	}

	// Run git status
	t.Log(w.Status())
	// Run git add .
	_, err = w.Add(".")
	if err != nil {
		t.Errorf("TestTwoParentsCommitDetectGoGit() add = %v", err)
	}
	// Run commit -m $message
	author := &object.Signature{Name: "test", Email: "test@test.com", When: time.Now()}
	firstHash, err := w.Commit("first commit", &git.CommitOptions{Author: author})
	if err != nil {
		t.Errorf("TestTwoParentsCommitDetectGoGit() commit = %v", err)
	}

	w, err = r.Worktree()
	if err != nil {
		t.Errorf("TestTwoParentsCommitDetectGoGit() worktree = %v", err)
	}

	file, err = fs.OpenFile("main.go", os.O_RDWR|os.O_CREATE, 0o755)
	if err != nil {
		t.Errorf("TestTwoParentsCommitDetectGoGit() open file main.go = %v", err)
	}
	err = file.Truncate(0)
	if err != nil {
		t.Errorf("TestTwoParentsCommitDetectGoGit() truncate main.go = %v", err)
	}
	_, err = file.Write([]byte(`package main\n\nimport "fmt"\nimport "log"\n\nfunc main() ` +
		`{\n	fmt.Println("Hello World")\n	log.Println("Hello World")\n}\n`))
	if err != nil {
		t.Errorf("TestTwoParentsCommitDetectGoGit() write file main.go = %v", err)
	}
	_ = file.Close()

	// Run git status
	t.Log(w.Status())
	// Run git add .
	if _, err = w.Add("."); err != nil {
		t.Errorf("TestTwoParentsCommitDetectGoGit() open file main.go = %v", err)
	}
	// Run commit -m $message
	editHash, err := w.Commit("edit main.go", &git.CommitOptions{Author: author})
	if err != nil {
		t.Errorf("TestTwoParentsCommitDetectGoGit() commit = %v", err)
	}

	w, err = r.Worktree()
	if err != nil {
		t.Errorf("TestTwoParentsCommitDetectGoGit() worktree = %v", err)
	}
	if _, err = w.Commit("merge", &git.CommitOptions{
		Parents: []plumbing.Hash{firstHash, editHash},
		Author:  author,
	}); err != nil {
		t.Errorf("TestTwoParentsCommitDetectGoGit() merge = %v", err)
	}

	// create the gitModel
	gitModel, err := local.NewGitModel(r)
	if err != nil {
		t.Errorf("TestTwoParentsCommitDetectGoGit() model = %v", err)
	}

	enrichedModel := enriched.NewEnrichedModel(*gitModel, remote.RemoteModel{})

	detector := NewCommitDetector(BranchCommitDetect())
	err = detector.Run(enrichedModel)
	if err != nil {
		t.Errorf("TestTwoParentsCommitDetectGoGit() run = %v", err)
	}

	log.Println(detector.Result())
}

func TestDiffMatchesMessageDetect(t *testing.T) {
	r := utils.FetchRepository(t, "https://github.com/Git-Gopher/tests", "test/commit-message-matches-diff/0")
	tests := []struct {
		name string
		want CommitDetect
	}{
		{"Commit message diff", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the gitModel
			gitModel, err := local.NewGitModel(r)
			if err != nil {
				t.Errorf(" TestDiffMatchesMessageDetect() create model = %v", err)
			}

			enrichedModel := enriched.NewEnrichedModel(*gitModel, remote.RemoteModel{})

			detector := NewCommitDetector(DiffMatchesMessageDetect())
			if err = detector.Run(enrichedModel); err != nil {
				t.Errorf(" TestDiffMatchesMessageDetect() run detector = %v", err)
			}

			t.Log(detector.Result())
		})
	}
}

func TestBinaryDetect(t *testing.T) {
	r := utils.FetchRepository(t, "https://github.com/Git-Gopher/tests", "test/committed-binary/1")
	tests := []struct {
		name string
		want CommitDetect
	}{
		{"BinaryDetect", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the gitModel
			gitModel, err := local.NewGitModel(r)
			if err != nil {
				t.Errorf(" TestBinaryDetect() create model = %v", err)
			}

			enrichedModel := enriched.NewEnrichedModel(*gitModel, remote.RemoteModel{})

			detector := NewCommitDetector(BinaryDetect())
			if err = detector.Run(enrichedModel); err != nil {
				t.Errorf(" TestBinaryDetect() run detector = %v", err)
			}

			t.Log(detector.Result())
		})
	}
}
