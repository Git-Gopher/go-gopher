package detector

import (
	"time"

	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/violation"
)

type MockBranchModel struct {
	Ref           string
	Remote        string
	Hash          string
	CommitsBehind int       // Number of commits behind the primary branch
	LastChange    time.Time // Time of the head commit of the current branch
}

type BranchDetect func(branches MockBranchModel) (bool, violation.Violation, error)

// BranchDetector is used to run a detector on each branch metadata.
type BranchDetector struct {
	violated   int
	found      int
	total      int
	violations []violation.Violation

	detect BranchDetect
}

// TODO: We should change this to the enriched model.
func (b *BranchDetector) Run(model *local.GitModel) error {
	b.violated = 0
	b.found = 0
	b.total = 0
	b.violations = make([]violation.Violation, 0)

	return ErrNotImplemented
}

func (b *BranchDetector) Result() (int, int, int, []violation.Violation) {
	return b.violated, b.found, b.total, b.violations
}

func NewBranchDetector(detect BranchDetect) *BranchDetector {
	return &BranchDetector{
		violated:   0,
		found:      0,
		total:      0,
		violations: make([]violation.Violation, 0),
		detect:     detect,
	}
}
