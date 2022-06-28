package detector

import (
	"fmt"
	"log"
	"strings"

	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/violation"
)

// Branch name.
type BranchCompareDetect func(branches []local.Branch) (int, []violation.Violation, error)

// BranchCompareDetector is used to run a detector on multiple branch and compare each branch.
type BranchCompareDetector struct {
	violated   int
	found      int
	total      int
	violations []violation.Violation

	detect BranchCompareDetect
}

func (b *BranchCompareDetector) Run(model *enriched.EnrichedModel) error {
	b.violated = 0
	b.found = 0
	b.total = 0
	b.violations = make([]violation.Violation, 0)

	if model == nil {
		return nil
	}

	var err error
	b.found, b.violations, err = b.detect(model.Branches)
	if err != nil {
		return fmt.Errorf("Error BranchCompareDetector: %w", err)
	}

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

// Deprecated: this detect is a demo
// NewFeatureBranchNewDetect is used to detect if a branch has the prefix feature or feat.
func NewFeatureBranchNameDetect() BranchCompareDetect {
	return func(branches []local.Branch) (int, []violation.Violation, error) {
		branchRefs := []string{}
		featureNames := [...]string{"feature", "feat"}

	b:
		for _, branch := range branches {
			for _, featureName := range featureNames {
				if strings.Contains(branch.Name, featureName) {
					// contains featureNames part of branch
					continue b
				}
			}
			// does not contain featureNames
			branchRefs = append(branchRefs, branch.Name)
		}

		// TODO: report using warning (not violation)
		log.Println("branch without feature/feat:", branchRefs)

		return 0, nil, nil
	}
}
