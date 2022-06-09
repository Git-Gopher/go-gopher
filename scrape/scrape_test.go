package scrape

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func environment() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Error loading .env file")
	}

	if os.Getenv("GITHUB_TOKEN") == "" {
		log.Fatalln("Error loading env GITHUB_TOKEN")
	}
}

func TestScraper_ScrapeUsers(t *testing.T) {
	// XXX: This test can change from underneath us if we decide to edit it
	environment()
	s := NewScraper()
	users, err := s.ScrapeUsers()
	if err != nil {
		t.Errorf("Could not scrape user: %v", err)
	}

	if users != nil {
		return
	}

	t.FailNow()
}

func TestScraper_ScrapePullRequests(t *testing.T) {
	environment()
	s := NewScraper()
	prs, err := s.ScrapePullRequests("Git-Gopher", "tests")
	if err != nil {
		t.Errorf("Could not scrape pull requests: %v", err)
	}

	// https://github.com/Git-Gopher/tests/tree/test/linked-pull-request-issue/0
	for _, pr := range prs {
		if pr.Id == "PR_kwDOHePx_M45XMYa" {
			return
		}
	}

	t.FailNow()
}
