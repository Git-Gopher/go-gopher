package detector

import (
	"reflect"
	"testing"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
)

func Test_rankSimilar(t *testing.T) {
	type args struct {
		input  []string
		metric strutil.StringMetric
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{
			name: "example",
			args: args{
				input: []string{"main",
					"feature/browse-artwork",
					"feature/manage-artwork",
					"feature/test"},
				metric: metrics.NewLevenshtein(),
			},
			want: []float64{0, 0, 0},
		},
		{
			name: "without main",
			args: args{
				input: []string{
					"feature/browse-artwork",
					"feature/manage-artwork",
					"feature/test"},
				metric: metrics.NewLevenshtein(),
			},
			want: []float64{0, 0, 0},
		},
		{
			name: "our repo",
			args: args{
				input: []string{
					"caching-basic",
					"various-tidy",
					"feature/branch-detectors",
					"action-all-branches",
				},
				metric: metrics.NewLevenshtein(),
			},
			want: []float64{0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rankSimilar(tt.args.input, tt.args.metric); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("rankSimilar() = %v, want %v", got, tt.want)
			}
		})
	}
}
