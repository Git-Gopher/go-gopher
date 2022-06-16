package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/utils"
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

// TODO: Use enriched model instead
func NewCache(model *local.GitModel) *Cache {
	return &Cache{
		Commits: model.Commits,
	}
}

// Github action will load the cache into the file system, we point this function towards it.
func LoadCaches(location string) []Cache {
	return nil
}

// Writes the cache file to the root of the project.
func (c *Cache) Write() error {
	file, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return fmt.Errorf("Error marshaling cache: %w", err)
	}

	root, err := utils.Root()
	if err != nil {
		return fmt.Errorf("Error fetching project root path: %w", err)
	}

	cachePath := path.Join(root, cacheFile)
	if err := ioutil.WriteFile(cachePath, file, 0o600); err != nil {
		return fmt.Errorf("Error writing cache: %w", err)
	}

	return nil
}

func (c *Cache) Read() error {
	root, err := utils.Root()
	if err != nil {
		return fmt.Errorf("Error fetching project root path: %w", err)
	}

	cachePath := path.Join(root, cacheFile)
	file, err := os.Open(cachePath)
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
