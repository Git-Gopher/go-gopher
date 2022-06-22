package detector

import (
	"github.com/Git-Gopher/go-gopher/model/github"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/violation"
)

// Branch name.
type BranchCompareDetect func(branches []MockBranchCompareModel) (bool, violation.Violation, error)

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

// TODO: Github Workflow: Branches must have consistent names.
func NewBranchNameConsistencyDetect() BranchCompareDetect {
	return func(branches []MockBranchCompareModel) (bool, violation.Violation, error) {
		// TODO: Do some algorithm to see if the branch names are consistent enough.
		if len(branches) > 10 {
			// TODO: Author should be whoever created the branch
			return false, violation.NewCommonViolation("Branch message longer than 10", nil), nil
		}

		return true, nil, nil
	}
}
