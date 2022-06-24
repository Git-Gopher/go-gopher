package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	giturls "github.com/whilp/git-urls"
)

var ErrUnsupportedSchema = errors.New("unsupported schema")

// Load the environment variables from the .env file.
func Environment(location string) {
	if err := godotenv.Load(location); err != nil {
		log.Println("Error loading .env file")
	}

	if os.Getenv("GITHUB_TOKEN") == "" {
		log.Fatalln("Error loading env GITHUB_TOKEN")
	}
}

// Fetch the owner and the name from the given URL.
// Supports https and ssh URLs.
func OwnerNameFromUrl(rawUrl string) (string, string, error) {
	fmt.Printf("rawUrl: %v\n", rawUrl)
	var owner, name string

	url, err := giturls.Parse(rawUrl)
	if err != nil {
		return "", "", fmt.Errorf("Could not parse git URL: %w", err)
	}

	xs := strings.Split(url.Path, "/")
	fmt.Printf("xs: %v\n", xs)
	switch url.Scheme {
	case "ssh":
		owner = xs[0]
		name = xs[1][:len(xs[1])-4] // Remove ".git".

	case "https":
		name = xs[0]
		owner = xs[1]
	case "http":
		name = xs[0]
		owner = xs[1]
	default:
		return "", "", fmt.Errorf("%w: %v", ErrUnsupportedSchema, url.Scheme)
	}

	return owner, name, nil
}
