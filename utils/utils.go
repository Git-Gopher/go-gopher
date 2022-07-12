package utils

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/joho/godotenv"
	giturls "github.com/whilp/git-urls"
)

var (
	ErrUnsupportedSchema = errors.New("unsupported schema")
	ErrRepo              = errors.New("Repository is nil")
)

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

	// XXX: Hack to remove .git from url
	name = strings.ReplaceAll(name, ".git", "")

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
		Auth: &githttp.BasicAuth{
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

// Fetch the Url from the repository remotes.
// Returns the first remote.
func Url(repo *git.Repository) (string, error) {
	if repo == nil {
		return "", ErrRepo
	}
	remotes, err := repo.Remotes()
	if err != nil {
		return "", fmt.Errorf("Could not get git repository remotes: %w", err)
	}

	if len(remotes) == 0 {
		return "", fmt.Errorf("No remotes present: %w", err)
	}

	// Use the first remote, assuming it correct.
	urls := remotes[0].Config().URLs
	if len(urls) == 0 {
		return "", fmt.Errorf("No URLs present: %w", err)
	}

	return urls[0], nil
}

// Check if a given filepath exists.
func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, fmt.Errorf("Could not check file exists status: %w", err)
}

func DownloadFile(path string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed GET url: %w", err)
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("failed create file: %w", err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	return err
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("could not open file: %w", err)
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(filepath.Clean(path), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer func() {
				if err = f.Close(); err != nil {
					panic(err)
				}
			}()

			// nolint: gosec
			_, err = io.Copy(f, rc)
			if err != nil {
				return fmt.Errorf("failed copy bytes to file: %w", err)
			}
		}

		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
