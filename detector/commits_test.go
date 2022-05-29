package detector

import (
	"log"
	"os"
	"testing"

	"github.com/Git-Gopher/go-gopher/model"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/joho/godotenv"
)

func FetchRepository(t *testing.T, remote string) *git.Repository {
	t.Helper()

	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Error loading .env file")
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Errorf("Empty token")
	}

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: "non-empty",
			Password: token,
		},
		URL: remote,
	})
	if err != nil {
		t.Errorf("%v", err)
		t.FailNow()
	}

	return r
}

func TestTwoParentsCommitDetect(t *testing.T) {
	r := FetchRepository(t, "https://github.com/Git-Gopher/github-two-parents-merged")
	tests := []struct {
		name string
		want CommitDetect
	}{
		{"asdf", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create the model
			model, err := model.NewGitModel(r)
			if err != nil {
				t.Errorf("TestTwoParentsCommitDetect() create model = %v", err)
			}

			detector := NewCommitDetector(TwoParentsCommitDetect())
			if err = detector.Run(model); err != nil {
				t.Errorf("TestTwoParentsCommitDetect() run detector = %v", err)
			}

			t.Log(detector.Result())
		})
	}
}

// TODO: Split go-git things into a suite of test friendly functions.
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
	firstHash, err := w.Commit("first commit", &git.CommitOptions{})
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
	editHash, err := w.Commit("edit main.go", &git.CommitOptions{})
	if err != nil {
		t.Errorf("TestTwoParentsCommitDetectGoGit() commit = %v", err)
	}

	w, err = r.Worktree()
	if err != nil {
		t.Errorf("TestTwoParentsCommitDetectGoGit() worktree = %v", err)
	}
	if _, err = w.Commit("merge", &git.CommitOptions{
		Parents: []plumbing.Hash{firstHash, editHash},
	}); err != nil {
		t.Errorf("TestTwoParentsCommitDetectGoGit() merge = %v", err)
	}

	// create the model
	model, err := model.NewGitModel(r)
	if err != nil {
		t.Errorf("TestTwoParentsCommitDetectGoGit() model = %v", err)
	}

	detector := NewCommitDetector(TwoParentsCommitDetect())
	err = detector.Run(model)
	if err != nil {
		t.Errorf("TestTwoParentsCommitDetectGoGit() run = %v", err)
	}

	log.Println(detector.Result())
}
