package model

import (
	"fmt"
	"testing"

	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/utils"
)

func TestPopulateAuthors(t *testing.T) {
	repoOwner := "Git-Gopher"
	repoName := "go-gopher"
	r := utils.FetchRepository(t, fmt.Sprintf("https://github.com/%s/%s", repoOwner, repoName), "main")

	enrichedModel, err := FetchEnrichedModel(r, repoOwner, repoName)
	if err != nil {
		t.Errorf("TestPopulateAuthors() fetch enriched model = %v", err)
	}

	authors := enriched.PopulateAuthors(enrichedModel)
	if authors == nil {
		t.Errorf("TestPopulateAuthors() = %v", authors)
	}

	t.Logf("authors: %+v", authors)
}
