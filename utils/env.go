package utils

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

//go:embed .env.reference
var envByte []byte

// Load the environment variables from the .env file.
func Environment(location string) {
	if err := Env(location); err != nil {
		log.Println("Error loading .env file")
	}

	if os.Getenv("GITHUB_TOKEN") == "" {
		log.Fatalln("Error loading env GITHUB_TOKEN")
	}
}

// Load an optional environment variables from the .env file.
func Env(location string) error {
	if err := godotenv.Load(location); err != nil {
		return fmt.Errorf("Error loading %s file: %w", location, err)
	}

	return nil
}

// Generate a default .env file.
func GenerateEnv(location string) error {
	if _, err := os.Stat(location); errors.Is(err, os.ErrNotExist) {
		dir := filepath.Dir(location)
		parent := filepath.Base(dir)

		if err2 := os.MkdirAll(parent, os.ModePerm); err2 != nil {
			return fmt.Errorf("can't create options dir: %w", err2)
		}
	} else if err != nil {
		overwrite := Confirm(fmt.Sprintf("Env: %s already exists. Overwrite?", location), 2)
		if !overwrite {
			return ErrSkipped
		}
	}

	if location == "" {
		location = ".env"
	}

	err := os.WriteFile(location, envByte, 0o600)
	if err != nil {
		return fmt.Errorf("can't write default options: %w", err)
	}

	return nil
}
