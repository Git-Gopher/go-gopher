package utils

import (
	"testing"
)

func TestOwnerNameFromURL(t *testing.T) {
	owner, name, err := OwnerNameFromUrl("https://github.com/Git-Gopher/go-gopher")
	if err != nil {
		t.Error(err)
	}
	if owner != "Git-Gopher" && name != "go-gopher" {
		t.Errorf("OwnerName() = %v, %v, want %v, %v", owner, name, "Git-Gopher", "go-gopher")
	}

	owner, name, err = OwnerNameFromUrl("git@github.com:Git-Gopher/go-gopher.git")
	if err != nil {
		t.Error(err)
	}
	if owner != "Git-Gopher" && name != "go-gopher" {
		t.Errorf("OwnerName() = %v, %v, want %v, %v", owner, name, "Git-Gopher", "go-gopher")
	}
}

func TestContains(t *testing.T) {
	type args struct {
		s  string
		xs []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"TestContainsJar", args{"thing.jar", []string{".jar", ".exe"}}, true},
		{"TestContainsExe", args{"thing.jar", []string{".jar", ".exe"}}, true},
		{"TestContainsNone", args{"thing", []string{".jar", ".exe"}}, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.s, tt.args.xs); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
