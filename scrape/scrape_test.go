package scrape

import (
	"testing"
)

func TestScraper_ScrapeUsers(t *testing.T) {
	// XXX: This test can change from underneath us if we decide to edit it
	s := NewScraper("https://github.com/Git-Gopher/github-two-parents-merged")
	users, err := s.ScrapeUsers()

	if err != nil {
		t.Errorf("Could not scrape user: %v", err)
	}

	t.Log(users)

}
