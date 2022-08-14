package utils

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"math"
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
	log "github.com/sirupsen/logrus"
	giturls "github.com/whilp/git-urls"
)

var (
	GoogleFormURL        = "https://forms.gle/gx8P86PefbBPH9J88"
	ErrUnsupportedSchema = errors.New("unsupported schema")
	ErrIllegalPath       = errors.New("illegal path")
	ErrRepo              = errors.New("repository is nil")
)

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
		SingleBranch:  true,
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
	// nolint: gosec
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed GET url: %w", err)
	}
	//nolint: errcheck
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("failed create file: %w", err)
	}

	//nolint: errcheck, gosec
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy body buffer: %w", err)
	}

	return nil
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer func() {
		if err = r.Close(); err != nil {
			panic(err)
		}
	}()

	err = os.MkdirAll(dest, 0o750)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("could not open file: %w", err)
		}
		defer func() {
			if err = rc.Close(); err != nil {
				panic(err)
			}
		}()

		// nolint: gosec
		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%w: %v", ErrIllegalPath, path)
		}

		// nolint: nestif
		if f.FileInfo().IsDir() {
			err = os.MkdirAll(path, f.Mode())
			if err != nil {
				return fmt.Errorf("failed to create directories: %w", err)
			}
		} else {
			err = os.MkdirAll(filepath.Dir(path), f.Mode())
			if err != nil {
				return fmt.Errorf("failed to create directories: %w", err)
			}
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

func RoundTime(input float64) int {
	var result float64

	if input < 0 {
		result = math.Ceil(input - 0.5)
	} else {
		result = math.Floor(input + 0.5)
	}

	i, _ := math.Modf(result)

	return int(i)
}

// Get the env variable for github workspace or return the default value.
func EnvGithubWorkspace() string {
	workspace := os.Getenv("GITHUB_WORKSPACE")
	if workspace == "" {
		workspace = "./"
		log.Println("GITHUB_WORKSPACE is not set, falling back to current directory...")
	}

	return workspace
}

// Contains for array of substrs.
func Contains(s string, xs []string) bool {
	for _, x := range xs {
		if strings.Contains(s, x) {
			return true
		}
	}

	return false
}
