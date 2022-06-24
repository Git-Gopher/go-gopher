package github

import (
	"fmt"
	"log"
	"testing"

	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/go-git/go-git/v5"
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

func TestThing(t *testing.T) {
	// Repo
	repo, err := git.PlainOpen("../../")
	if err != nil {
		log.Fatalf("cannot read repo: %v\n", err)
	}

	config, err := repo.Config()
	fmt.Printf("config.URLs: %v\n", config.URLs)
	remotes, err := repo.Remotes()
	if err != nil {
		log.Fatalf("Could not get git repository remotes: %v\n", err)
	}

	if len(remotes) == 0 {
		log.Fatalf("No remotes present: %v\n", err)
	}

	// XXX: Use the first remote, assuming origin.
	urls := remotes[0].Config().URLs
	if len(urls) == 0 {
		log.Fatalf("No URLs present: %v\n", err)
	}

	url := urls[0]

	utils.OwnerNameFromUrl(url)
	// gitModel, err := local.NewGitModel(repo)
	// if err != nil {
	// 	log.Fatalf("Could not create GitModel: %v\n", err)
	// }

	// fmt.Printf("gitModel: %v\n", gitModel)
}
