package scrape

import (
	"reflect"
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

func TestScraper_ScrapePullRequests(t *testing.T) {
	s := NewScraper("https://github.com/Git-Gopher/tests")
	prs, err := s.ScrapePullRequests("Git-Gopher", "tests")
	if err != nil {
		t.Errorf("Could not scrape pull requests: %v", err)
	}

	// https://github.com/Git-Gopher/tests/tree/test/linked-pull-request-issue/0
	for _, pr := range prs {
		if (reflect.DeepEqual(pr, pullRequest{
			Id:                      "PR_kwDOHePx_M45XMYa",
			Title:                   "test/linked-pull-request-issue/modify",
			Body:                    "closes: #1",
			ClosingIssuesReferences: pr.ClosingIssuesReferences,
		})) {
			return
		}
	}

	t.FailNow()
}
