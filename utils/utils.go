package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/joho/godotenv"
	giturls "github.com/whilp/git-urls"
)

var ErrUnsupportedSchema = errors.New("unsupported schema")

// Load the environment variables from the .env file.
func Environment(location string) {
	if err := godotenv.Load(location); err != nil {
		log.Println("Error loading .env file")
	}

	if os.Getenv("GITHUB_TOKEN") == "" {
		log.Fatalln("Error loading env GITHUB_TOKEN")
	}
}

// Fetch the owner and the name from the given URL.
// Supports https and ssh URLs.
func OwnerNameFromUrl(rawUrl string) (string, string, error) {
	var owner, name string

	fmt.Printf("rawUrl: %v\n", rawUrl)
	url, err := giturls.Parse(rawUrl)
	if err != nil {
		return "", "", fmt.Errorf("Could not parse git URL: %w", err)
	}

	xs := strings.Split(url.Path, "/")
	switch url.Scheme {
	case "ssh":
		owner = xs[0]
		name = xs[1][:len(xs[1])-4] // Remove ".git".
	case "https":
		owner = xs[1]
		name = xs[2]
	case "http":
		owner = xs[1]
		name = xs[2]
	default:
		return "", "", fmt.Errorf("%w: %v", ErrUnsupportedSchema, url.Scheme)
	}

	return owner, name, nil
}

func FetchRepository(t *testing.T, remote, branch string) *git.Repository {
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
		URL:           remote,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
	})
	if err != nil {
		t.Errorf("%v", err)
		t.FailNow()
	}

	return r
}
