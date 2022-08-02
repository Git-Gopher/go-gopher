package utils

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"

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
	if location == "" {
		location = ".env"
	}

	err := ioutil.WriteFile(location, envByte, 0o600)
	if err != nil {
		return fmt.Errorf("can't write default options: %w", err)
	}

	return nil
}
