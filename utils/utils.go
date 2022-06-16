package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
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

// Find the root of the project.
func Root() (string, error) {
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("Failed to get root directory")
	}

	return filepath.Join(filepath.Dir(b), "../"), nil
}
