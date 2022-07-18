package assess

import (
	"github.com/Git-Gopher/go-gopher/detector"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/utils"
	"github.com/Git-Gopher/go-gopher/violation"
)

type Mark struct {
	Username   string
	Violations []violation.Violation
	Total      int
}

type MarkerCtx struct {
	Model        *enriched.EnrichedModel
	Contribution *contribution
	Author       utils.Authors
}

type Marker func(MarkerCtx) (string, []Mark)

// List all markers.
var _ []Marker = []Marker{
	Atomicity,          // E-A1 Diff Distance/Atomicity.
	DescriptiveCommit,  // E-A1 Commit messages are descriptive.
	GeneratedFiles,     // E-A1 Generated files should not be committed.
	RegularBranchNames, // E-A2 Regular branch names.

	FeatureBranching,  // W-A2 Feature branching.
	PullRequestReview, // W-A2 Pull request review.
}

// DetectorMarker is a helper method to run a detector as a marker.
func DetectorMarker(
	m MarkerCtx, // ctx
	d detector.Detector, // detector
	c map[string]int, // contribution map
) []Mark {
	if err := d.Run(m.Model); err != nil {
		return nil
	}

	// violations map to author
	violationsMap := make(map[string][]violation.Violation)

	_, _, _, violations := d.Result() // nolint:dogsled
	for _, v := range violations {
		username, err := m.Author.Find(v.Email())
		if err != nil {
			continue
		}

		violationsMap[*username] = append(violationsMap[*username], v)
	}

	marks := make([]Mark, 0, len(violationsMap))
	for username, violations := range violationsMap {
		emails, err := m.Author.Details(username)
		if err != nil {
			continue
		}

		count := 0 // total contributions for this author.
		for _, email := range emails {
			count += c[email]
		}

		marks = append(marks, Mark{
			Username:   username,
			Violations: violations,
			Total:      count,
		})
	}

	return marks
}
