package cache

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/utils"
)

const cacheFile = "cache.json"

// Cache belonging to a particular branch.
type Cache struct {
	Created time.Time           `json:"time"`
	Hashes  map[string]struct{} `json:"hashes"`
}

// Cache from a enriched model with empty missing hash set.
func NewCache(em *enriched.EnrichedModel) *Cache {
	hashes := make(map[string]struct{})
	for _, commit := range em.Commits {
		hashes[commit.Hash.HexString()] = struct{}{}
	}

	return &Cache{
		Created: time.Now(),
		Hashes:  hashes,
	}
}

// Github action will load the cache into the file system, we point this function towards it.
func Read() ([]*Cache, error) {
	loc := utils.EnvGithubWorkspace()
	path := path.Join(loc, cacheFile)

	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("error opening cache file: %w", err)
	}

	//nolint
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading cache file: %w", err)
	}

	var caches []*Cache
	if err = json.Unmarshal(b, &caches); err != nil {
		return nil, fmt.Errorf("error unmarshalling cache file: %w", err)
	}

	return caches, nil
}

func Write(caches []*Cache) error {
	loc := utils.EnvGithubWorkspace()
	path := path.Join(loc, cacheFile)

	bytes, err := json.MarshalIndent(caches, "", " ")
	if err != nil {
		return fmt.Errorf("error marshaling caches: %w", err)
	}

	if err := os.WriteFile(path, bytes, 0o600); err != nil {
		return fmt.Errorf("error writing cache: %w", err)
	}

	return nil
}
