package detector

import (
	"sort"
	"strings"

	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/violation"
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"gopkg.in/vmarkovtsev/go-lcss.v1"
)

// Branches must have consistent names.
// Research: https://stackoverflow.com/questions/29476737/similarities-in-strings-for-name-matching
// Methods: q-grams, longest common substring and longest common subsequence.
func NewBranchNameConsistencyDetect() BranchCompareDetect {
	return func(branches []local.Branch) (int, []violation.Violation, error) {
		branchRefs := make([]string, len(branches))
		for i, branch := range branches {
			branchRefs[i] = branch.Name
		}

		ranking := rankSimilar(branchRefs, metrics.NewLevenshtein())
		ranks := &sortByRanking{branchRefs, ranking}
		sort.Sort(ranks)
		ranked := ranks.name

		// half the top results
		ranked = ranked[:len(ranked)/2]

		substring := longestSubstring(ranked)
		if len(substring) <= 3 { // ignore if less than 3 char
			substring = ""
		}

		count := 0
		violations := []violation.Violation{}
		for i, branch := range branches {
			if substring != "" && strings.Contains(branch.Name, substring) {
				count++

				continue
			}
			// does not follow substring
			if ranking[i] < 0.175*float64(len(branches)) { // 0.175 is an adjustable value
				// not consistent with others.
				violations = append(
					violations,
					violation.NewBranchNameViolation(branch.Name, substring, branch.Head.Committer.Email),
				)
			}

			// TODO: warning not using substring
		}

		return count, violations, nil
	}
}

// sortByRanking sorts the input strings by their ranking.
// implements sort.Interface
// usage: sort.Sort(&sortByRanking{}).
type sortByRanking struct {
	name []string
	rank []float64
}

func (sbr sortByRanking) Len() int {
	return len(sbr.name)
}

func (sbr sortByRanking) Swap(i, j int) {
	sbr.name[i], sbr.name[j] = sbr.name[j], sbr.name[i]
	sbr.rank[i], sbr.rank[j] = sbr.rank[j], sbr.rank[i]
}

func (sbr sortByRanking) Less(i, j int) bool {
	return sbr.rank[i] < sbr.rank[j]
}

// rankSimilar ranks the similarity of the input strings.
func rankSimilar(input []string, metric strutil.StringMetric) []float64 {
	results := make([]float64, len(input))
	for i := 0; i < len(input); i++ {
		for j := i + 1; j < len(input); j++ {
			similarity := strutil.Similarity(input[i], input[j], metric)
			results[i] += similarity
			results[j] += similarity
		}
	}

	return results
}

// longestSubstring finds the longest common substring of the input strings.
func longestSubstring(input []string) string {
	b := make([][]byte, len(input))
	for i, str := range input {
		b[i] = []byte(str)
	}

	return string(lcss.LongestCommonSubstring(b...))
}
