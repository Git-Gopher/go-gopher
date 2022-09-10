package remote

import (
	"testing"

	"github.com/Git-Gopher/go-gopher/utils"
)

func TestScrapeRemoteModel(t *testing.T) {
	type args struct {
		owner string
		name  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"ScrapeRemoteModel",
			args{owner: "Git-Gopher", name: "tests"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utils.Environment("../../.env")
			_, err := ScrapeRemoteModel(tt.args.owner, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScrapeRemoteModel() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
		})
	}
}
