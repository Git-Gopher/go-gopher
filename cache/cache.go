package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
)

const cacheFile = "cache.json"

// caches: [
// 			{
// 					: []
// 			},
// 			{
// 					commits: []
// 			}
// 	]

type Cache struct {
	Commits []local.Commit `json:"commits"`
}

func NewCache(model *enriched.EnrichedModel) {
}

// Github action will load the cache into the file system, we point this function towards it.
func LoadCaches(location string) []Cache {
	return nil
}

// Will save to the repository root (where the application is run from).
func (c *Cache) Write() error {
	file, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return fmt.Errorf("Error marshaling cache: %w", err)
	}

	if err := ioutil.WriteFile(cacheFile, file, 0o600); err != nil {
		return fmt.Errorf("Error writing cache: %w", err)
	}

	return nil
}

func (c *Cache) Read() error {
	file, err := os.Open("cache.json")
	if err != nil {
		return fmt.Errorf("Error opening cache file: %w", err)
	}

	//nolint
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Error reading cache file: %w", err)
	}

	if err = json.Unmarshal(b, c); err != nil {
		return fmt.Errorf("Error unmarshalling cache file: %w", err)
	}

	return nil
}
