package cache

import (
	"testing"

	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestCache_Write(t *testing.T) {
	type fields struct {
		Commits []local.Commit
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"Write Basic", fields{Commits: []local.Commit{
			*local.NewCommit(
				&object.Commit{Message: "asdf"}),
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				Commits: tt.fields.Commits,
			}
			if err := c.Write(); (err != nil) != tt.wantErr {
				t.Errorf("Cache.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
