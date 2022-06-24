package utils

import (
	"testing"
)

func TestOwnerNameFromURL(t *testing.T) {
	owner, name, err := OwnerNameFromUrl("https://github.com/Git-Gopher/go-gopher")
	if err != nil {
		t.Error(err)
	}

	owner, name, err = OwnerNameFromUrl("git@github.com:Git-Gopher/go-gopher.git")
	if err != nil {
		t.Error(err)
	}
	if owner != "Git-Gopher" && name != "go-gopher" {
		t.Errorf("OwnerName() = %v, %v, want %v, %v", owner, name, "Git-Gopher", "go-gopher")
	}
}
