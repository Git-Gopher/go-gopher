package scrape

import (
	"testing"

	"github.com/Git-Gopher/go-gopher/utils"
)

func TestScraper_ScrapeUsers(t *testing.T) {
	// XXX: This test can change from underneath us if we decide to edit it
	utils.Environment("../.env")
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
	utils.Environment("../.env")
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
