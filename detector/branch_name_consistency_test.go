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
				input: []string{
					"main",
					"feature/browse-artwork",
					"feature/manage-artwork",
					"feature/test",
				},
				metric: metrics.NewLevenshtein(),
			},
			want: []float64{0.21969696969696972, 1.2727272727272727, 1.3181818181818183, 0.9924242424242425},
		},
		{
			name: "without main",
			args: args{
				input: []string{
					"feature/browse-artwork",
					"feature/manage-artwork",
					"feature/test",
				},
				metric: metrics.NewLevenshtein(),
			},
			want: []float64{1.2272727272727273, 1.2272727272727273, 0.9090909090909092},
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
			want: []float64{0.636302294197031, 0.4784075573549258, 0.5416666666666666, 0.6820175438596492},
		},
		{
			name: "game example",
			args: args{
				input: []string{
					"feat/sql-database",
					"feature/score-board",
					"feature/tile-animations",
					"feature/proof-of-concept",
				},
				metric: metrics.NewLevenshtein(),
			},
			want: []float64{0.9987604881769641, 1.208905415713196, 1.2010869565217392, 1.125},
		},
		{
			name: "game example",
			args: args{
				input: []string{
					"feat/sql-database",
					"feature/score-board",
					"feature/tile-animations",
					"feature/proof-of-concept",
					"feature/user-and-leaderboard-controller",
					"feature/browse-artwork",
					"feature/manage-artwork",
				},
				metric: metrics.NewLevenshtein(),
			},
			want: []float64{
				1.8519073413238172,
				2.638975345783126,
				2.516583054626533,
				2.342948717948718,
				2.205128205128205,
				2.886477652782001,
				2.8089718252761733,
			},
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

func Test_longestSubstring(t *testing.T) {
	type args struct {
		input []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"features", args{[]string{"feature/score-board", "feature/tile-animations"}}, "feature/"},
		{"nothing", args{[]string{"caching-basic", "various-tidy"}}, "-"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := longestSubstring(tt.args.input); got != tt.want {
				t.Errorf("longestSubstring() = %v, want %v", got, tt.want)
			}
		})
	}
}
