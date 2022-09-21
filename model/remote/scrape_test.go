package remote

import (
	"context"
	"testing"

	"github.com/Git-Gopher/go-gopher/utils"
)

func TestScraper_FetchPullRequests(t *testing.T) {
	utils.Environment("../../.env")
	s := NewScraper()
	prs, err := s.FetchPullRequests(context.TODO(), "subquery", "subql")
	if err != nil {
		t.Error(err)
	}
	t.Logf("prs: %v\n", prs)
}

func TestScraper_FetchCommitters(t *testing.T) {
	utils.Environment("../../.env")
	s := NewScraper()
	committers, err := s.FetchCommitters(context.TODO(), "Git-Gopher", "go-gopher")
	if err != nil {
		t.Error(err)
	}
	t.Logf("committers: %v\n", committers)
}

func TestScraper_FetchPopularRepositories(t *testing.T) {
	utils.Environment("../../.env")
	s := NewScraper()
	repos, err := s.FetchPopularRepositories(context.TODO(), 100, 1, 1, 10)
	if err != nil {
		t.Error(err)
	}

	t.Logf("repos: %v\n", repos)
}
