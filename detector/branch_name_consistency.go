package detector

import (
	"sort"
	"strings"

	"github.com/Git-Gopher/go-gopher/violation"
	"github.com/adrg/strutil/metrics"
)

// Branches must have consistent names.
// Research: https://stackoverflow.com/questions/29476737/similarities-in-strings-for-name-matching
// Methods: q-grams, longest common substring and longest common subsequence.
func NewBranchNameConsistencyDetect() BranchCompareDetect {
	return func(branches []MockBranchCompareModel) (int, []violation.Violation, error) {
		branchRefs := make([]string, len(branches))
		for i, branch := range branches {
			branchRefs[i] = branch.Ref
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
			if substring != "" && strings.Contains(branch.Ref, substring) {

				count++
				continue
			}
			// does not follow substring
			if ranking[i] < 0.175*float64(len(branches)) { // 0.175 is an adjustable value
				// not consistent with others
				violations = append(violations, violation.NewBranchNameViolation(branch.Ref, substring))
			}

			// TODO: warning not using substring
		}

		return count, violations, nil
	}
}

type ranker struct {
	name []string
	rank []float64
}

type sortByRanking ranker

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
