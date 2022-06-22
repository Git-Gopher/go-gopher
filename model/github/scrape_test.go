package github

import (
	"fmt"
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
	fmt.Printf("prs: %v\n", prs)
}
