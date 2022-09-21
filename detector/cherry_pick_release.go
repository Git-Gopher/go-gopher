package detector

import (
	"errors"

	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/violation"
)

var (
	ErrCherryPickReleaseModelNil        = errors.New("cherry pick model is nil")
	ErrCherryPickReleaseNoReleaseBranch = errors.New("cherry pick detector: no release branch")
	ErrCherryPickReleaseNoDevBranch     = errors.New("cherry pick detector: no dev branch")
)

// CherryPickReleaseDetector is a detector that detects cherry pick across release and dev branch.
// found / total = number of cherry picked commits found across release and dev.
// CherryPickReleaseDetector rely on patch-id which does not detect cherry pick if there is a conflict.
type CherryPickReleaseDetector struct {
	name       string
	violated   int // non feature branches aka develop/release etc. (does not account default branch)
	found      int // total feature branches
	total      int // total branches
	violations []violation.Violation
}

// NewCherryPickReleaseDetector creates a new cherry pick detector.
// This detector compares release and dev branch and does not rely on the
// `detector.NewDetector(detector.Detect)` pattern.
func NewCherryPickReleaseDetector(name string) *CherryPickReleaseDetector {
	return &CherryPickReleaseDetector{
		name:       name,
		violated:   0,
		found:      0,
		total:      0,
		violations: make([]violation.Violation, 0),
	}
}

//nolint:gocognit // VERY COMPLEX
func (cp *CherryPickReleaseDetector) Run(em *enriched.EnrichedModel) error {
	if em == nil {
		return ErrCherryPickReleaseModelNil
	}

	if em.ReleaseGraph == nil {
		return ErrCherryPickReleaseNoReleaseBranch
	}

	if em.MainGraph == nil {
		return ErrCherryPickReleaseNoDevBranch
	}

	cp.violated = 0
	cp.found = 0
	cp.total = 0
	cp.violations = make([]violation.Violation, 0)

	patchIDMap := make(map[string]local.Commit)
	matched := make(map[string][]local.Commit)

	// generate matched map with all cherry picked commits
	for _, commit := range em.Commits {
		if commit.PatchID != nil {
			pid := *commit.PatchID
			if c, ok := patchIDMap[pid]; ok {
				// duplicate patch id means cherry pick?
				if _, ok := matched[pid]; !ok {
					matched[pid] = []local.Commit{c}
				}
				matched[pid] = append(matched[pid], commit)
			}

			patchIDMap[pid] = commit
		}
	}

	// generate commitMap which maps commit hash to pid
	commitMap := make(map[string]string)

	for pid, commits := range matched {
		for _, commit := range commits {
			commitMap[commit.Hash.HexString()] = pid
		}
	}

	// cherry-picked checklist
	cherryPicked := make(map[string]struct{})

	// for loop through all commits in main graph
	mainCommits := make(map[string]*local.CommitGraph)
	next := []*local.CommitGraph{}
	for _, commit := range em.MainGraph.Head.ParentCommits {
		mainCommits[commit.Hash] = commit
		if pid, ok := commitMap[commit.Hash]; ok {
			cherryPicked[pid] = struct{}{}
		}

		next = append(next, commit)
	}

	for len(next) != 0 {
		commit := next[0]
		next = next[1:]

		if _, ok := mainCommits[commit.Hash]; ok {
			continue
		}

		mainCommits[commit.Hash] = commit
		if pid, ok := commitMap[commit.Hash]; ok {
			cherryPicked[pid] = struct{}{}
		}

		next = append(next, commit.ParentCommits...)
	}

	// for loop through all commits in release graph
	releaseCommits := make(map[string]*local.CommitGraph)
	next = []*local.CommitGraph{}
	for _, commit := range em.ReleaseGraph.Head.ParentCommits {
		releaseCommits[commit.Hash] = commit
		if pid, ok := commitMap[commit.Hash]; ok {
			if _, ok := cherryPicked[pid]; ok {
				cp.found++ // found cherry pick +1
			}
		}

		cp.total++ // count commits in release branch
		next = append(next, commit)
	}

	for len(next) != 0 {
		commit := next[0]
		next = next[1:]

		if _, ok := releaseCommits[commit.Hash]; ok {
			continue
		}

		releaseCommits[commit.Hash] = commit
		if pid, ok := commitMap[commit.Hash]; ok {
			if _, ok := cherryPicked[pid]; ok {
				cp.found++ // found cherry pick +1
			}
		}

		cp.total++ // count commits in release branch
		next = append(next, commit.ParentCommits...)
	}

	return nil
}

func (cp *CherryPickReleaseDetector) Result() (int, int, int, []violation.Violation) {
	return cp.violated, cp.found, cp.total, cp.violations
}

func (cp *CherryPickReleaseDetector) Name() string {
	return cp.name
}
