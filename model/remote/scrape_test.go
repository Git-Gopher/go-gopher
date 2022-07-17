package remote

import (
	"testing"

	"github.com/Git-Gopher/go-gopher/utils"
)

func TestScraper_FetchPullRequests(t *testing.T) {
	utils.Environment("../../.env")
	s := NewScraper()
	prs, err := s.FetchPullRequests("subquery", "subql")
	if err != nil {
		t.Error(err)
	}
	t.Logf("prs: %v\n", prs)
}

func TestScraper_FetchCommitters(t *testing.T) {
	utils.Environment("../../.env")
	s := NewScraper()
	committers, err := s.FetchCommitters("Git-Gopher", "go-gopher")
	if err != nil {
		t.Error(err)
	}
	t.Logf("committers: %v\n", committers)
}
