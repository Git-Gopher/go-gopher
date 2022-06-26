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
)

const cacheFile = "cache.json"

// Cache is a single cache file containing a set of caches.
type Cache struct {
	Hashes  []local.Hash `json:"hashes"`
	Created time.Time    `json:"time"`
}

func NewCache(model *enriched.EnrichedModel) *Cache {
	hashes := make([]local.Hash, len(model.Commits))
	for i, c := range model.Commits {
		hashes[i] = c.Hash
	}

	return &Cache{
		Hashes:  hashes,
		Created: time.Now(),
	}
}

// Github action will load the cache into the file system, we point this function towards it.
func ReadCaches() ([]*Cache, error) {
	loc := Location()
	cachePath := path.Join(loc, cacheFile)
	file, err := os.Open(filepath.Clean(cachePath))
	if err != nil {
		return nil, fmt.Errorf("Error opening cache file: %w", err)
	}

	//nolint
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Error reading cache file: %w", err)
	}

	var caches []*Cache
	if err = json.Unmarshal(b, &caches); err != nil {
		return nil, fmt.Errorf("Error unmarshalling cache file: %w", err)
	}

	return caches, nil
}

func WriteCaches(caches []*Cache) error {
	bytes, err := json.MarshalIndent(caches, "", " ")
	if err != nil {
		return fmt.Errorf("Error marshaling caches: %w", err)
	}

	loc := Location()
	cachePath := path.Join(loc, cacheFile)
	if err := ioutil.WriteFile(cachePath, bytes, 0o600); err != nil {
		return fmt.Errorf("Error writing cache: %w", err)
	}

	return nil
}

// Attempt to find the location for the cache file.
// $GITHUB_WORKSPACE, otherwise fallback onto $CACHE_PATH.
func Location() string {
	workspace := os.Getenv("GITHUB_WORKSPACE")
	if workspace == "" {
		workspace = os.Getenv("CACHE_PATH")
	}

	return workspace
}
