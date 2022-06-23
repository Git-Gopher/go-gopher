package detector

import (
	"log"
	"strings"

	"github.com/Git-Gopher/go-gopher/model/github"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/violation"
	"github.com/adrg/strutil"
	"gopkg.in/vmarkovtsev/go-lcss.v1"
)

// Branch name.
type BranchCompareDetect func(branches []MockBranchCompareModel) (int, []violation.Violation, error)

// XXX: Move into either local or github model.
type MockBranchCompareModel struct {
	Ref    string
	Remote string
	Hash   string
}

// BranchCompareDetector is used to run a detector on multiple branch and compare each branch.
type BranchCompareDetector struct {
	violated   int
	found      int
	total      int
	violations []violation.Violation

	detect BranchCompareDetect
}

// TODO: We should change this to the enriched model.
func (b *BranchCompareDetector) Run(model *local.GitModel) error {
	b.violated = 0
	b.found = 0
	b.total = 0
	b.violations = make([]violation.Violation, 0)

	return ErrNotImplemented
}

func (b *BranchCompareDetect) Run2(ghm *github.GithubModel) error {
	return nil
}

func (b *BranchCompareDetector) Result() (int, int, int, []violation.Violation) {
	return b.violated, b.found, b.total, b.violations
}

func NewBranchCompareDetector(detect BranchCompareDetect) *BranchCompareDetector {
	return &BranchCompareDetector{
		violated:   0,
		found:      0,
		total:      0,
		violations: make([]violation.Violation, 0),
		detect:     detect,
	}
}

func NewFeatureBranchNameDetect() BranchCompareDetect {
	return func(branches []MockBranchCompareModel) (int, []violation.Violation, error) {
		branchRefs := []string{}
		featureNames := [...]string{"feature", "feat"}

	b:
		for _, branch := range branches {
			for _, featureName := range featureNames {
				if strings.Contains(branch.Ref, featureName) {
					// contains featureNames part of branch
					continue b
				}
			}
			// does not contain featureNames
			branchRefs = append(branchRefs, branch.Ref)
		}

		// TODO: report using warning (not violation)
		log.Println("branch without feature/feat:", branchRefs)

		return 0, nil, nil
	}
}

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

func longestSubstring(input []string) string {
	b := make([][]byte, len(input))
	for i, str := range input {
		b[i] = []byte(str)
	}

	return string(lcss.LongestCommonSubstring(b...))
}
