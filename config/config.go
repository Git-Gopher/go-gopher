package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	DefaultConfigPath = "config.yml"

	//go:embed config.yml
	defaultConfig []byte
)

// Config for enabling and disabling detectors.
type Config struct {
	Detectors map[string]struct {
		Enabled bool
		Weight  int
	}
}

func Read(path string) (*Config, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("Failed to read config file: %v", err)
	}

	t := Config{}

	if err = yaml.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("Failed to parse config file: %v", err)
	}

	return &t, nil
}

// Read a default config from embedded default config.
func Default() (*Config, error) {
	t := Config{}
	if err := yaml.Unmarshal(defaultConfig, &t); err != nil {
		return nil, fmt.Errorf("Failed to parse embedded config file: %w", err)
	}

	return &t, nil
}
