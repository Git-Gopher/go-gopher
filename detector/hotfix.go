package detector

import (
	"errors"
	"sort"
	"strings"

	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/violation"
)

var ErrHotfixModelNil = errors.New("hotfix model is nil")

// HotfixDetector that finds number of hotfixes using tag number.
type HotfixDetector struct {
	name       string
	violated   int // non feature branches aka develop/release etc. (does not account default branch)
	found      int // total feature branches
	total      int // total branches
	violations []violation.Violation
}

// NewHotfixDetector creates a new hotfix detector.
func NewHotfixDetector(name string) *HotfixDetector {
	return &HotfixDetector{
		name:       name,
		violated:   0,
		found:      0,
		total:      0,
		violations: make([]violation.Violation, 0),
	}
}

func (cp *HotfixDetector) Run(em *enriched.EnrichedModel) error {
	if em == nil {
		return ErrHotfixModelNil
	}

	sortedTags := em.Tags

	sort.Slice(sortedTags, func(i, j int) bool {
		return sortedTags[i].Name < sortedTags[j].Name
	})

	for i, tag := range sortedTags {
		if !strings.HasPrefix(tag.Name, "v") {
			// skip tags with no v prefix
			// must follow v1.0.0 format
			continue
		}

		// total number of tags
		cp.total++

		if i == 0 {
			// first tag cannot be a hotfix
			continue
		}

		prev := sortedTags[i-1].Head.Committer
		curr := tag.Head.Committer

		if prev.When.After(curr.When) {
			// The previous tag is older than the next tag
			// this means the prev tag is a hotfix
			cp.found++
		}
	}

	return nil
}

func (cp *HotfixDetector) Result() (int, int, int, []violation.Violation) {
	return cp.violated, cp.found, cp.total, cp.violations
}

func (cp *HotfixDetector) Name() string {
	return cp.name
}
