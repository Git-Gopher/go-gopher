package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/utils"
)

const cacheFile = "cache.json"

// Cache belonging to a particular branch.
type Cache struct {
	Created       time.Time               `json:"time"`
	Hashes        map[local.Hash]struct{} `json:"hashes"`
	MissingHashes map[local.Hash]struct{} `json:"missing_hashes"`
}

// Cache from a enriched model with empty missing hash set.
func NewCache(em *enriched.EnrichedModel) *Cache {
	hashes := make(map[local.Hash]struct{})
	for _, commit := range em.Commits {
		hashes[commit.Hash] = struct{}{}
	}

	return &Cache{
		Created: time.Now(),
		Hashes:  hashes,
	}
}

// Github action will load the cache into the file system, we point this function towards it.
func Read() (*Cache, error) {
	loc := utils.EnvGithubWorkspace()
	path := path.Join(loc, cacheFile)
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("error opening cache file: %w", err)
	}

	//nolint
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading cache file: %w", err)
	}

	var cache *Cache
	if err = json.Unmarshal(b, cache); err != nil {
		return nil, fmt.Errorf("error unmarshalling cache file: %w", err)
	}

	return cache, nil
}

// Update the cache by comparing the old cache to current model.
func (c *Cache) Update(em *enriched.EnrichedModel) {
	for _, commit := range em.Commits {
		if _, ok := c.Hashes[commit.Hash]; ok {
			delete(c.MissingHashes, commit.Hash)
		} else {
			c.MissingHashes[commit.Hash] = struct{}{}
		}
	}

	hashes := make(map[local.Hash]struct{})
	for _, commit := range em.Commits {
		hashes[commit.Hash] = struct{}{}
	}

	c.Hashes = hashes
	c.Created = time.Now()
}

func (c *Cache) Write() error {
	loc := utils.EnvGithubWorkspace()

	bytes, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return fmt.Errorf("Error marshaling caches: %w", err)
	}

	path := path.Join(loc, cacheFile)
	if err := ioutil.WriteFile(path, bytes, 0o600); err != nil {
		return fmt.Errorf("Error writing cache: %w", err)
	}

	return nil
}
