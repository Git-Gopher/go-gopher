package detector

import (
	"sort"
	"testing"

	"github.com/Git-Gopher/go-gopher/model/local"
)

func TestHotfix(t *testing.T) {
	sortedTags := []local.Tag{
		{
			Name: "v2.0.0",
		}, {
			Name: "v1.0.0",
		}, {
			Name: "v1.0.5",
		}, {
			Name: "v1.0.2",
		},
	}

	sort.Slice(sortedTags, func(i, j int) bool {
		return sortedTags[i].Name < sortedTags[j].Name
	})

	for _, tag := range sortedTags {
		t.Logf("%s\n", tag.Name)
	}
}
