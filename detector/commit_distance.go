package detector

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/violation"
)

type CommitDistanceCalculator func(commit *local.Commit) (distance float64, err error)

type CommitDistanceDetector struct {
	violated   int
	found      int
	total      int
	violations []violation.Violation

	distance []float64

	detect CommitDistanceCalculator
}

func (cd *CommitDistanceDetector) Run(model *enriched.EnrichedModel) error {
	// Struct should be reset before each run, incase we are running it with a different model.
	cd.violated = 0
	cd.found = 0
	cd.total = 0
	cd.violations = make([]violation.Violation, 0)

	cd.distance = make([]float64, len(model.Commits))

	for i, c := range model.Commits {
		c := c
		distance, err := cd.detect(&c)
		cd.total++
		if err != nil {
			return err
		}

		cd.distance[i] = distance
	}

	// find outliers from the distance using.
	// fmt.Println("++++++++++++++++++++++++++++++++++++++++")
	// fmt.Println(cd.distance)

	return nil
}

func (cd *CommitDistanceDetector) Result() (int, int, int, []violation.Violation) {
	return cd.violated, cd.found, cd.total, cd.violations
}

func NewCommitDistanceDetector(detect CommitDistanceCalculator) *CommitDistanceDetector {
	return &CommitDistanceDetector{
		violated:   0,
		found:      0,
		total:      0,
		violations: make([]violation.Violation, 0),
		detect:     detect,
	}
}

// 6th methods: use distance between diff to find an average.
// nolint:gocognit // this function is complex
func DiffDistanceCalculation() CommitDistanceCalculator {
	return func(commit *local.Commit) (distance float64, err error) {
		if commit.DiffToParents == nil {
			// no diff
			return 0.0, nil
		}

		averages := map[string]float64{}

		for _, diff := range commit.DiffToParents {
			filename := diff.Name

			var max int64
			var min int64
			for _, point := range diff.Points {
				value := point.NewPosition

				if max < value {
					max = value
				}
				if min > value {
					min = value
				}
			}

			average := float64(max-min) / float64(len(diff.Points))

			// clamp: average cannot be zero for calculations
			if average == 0 {
				average = 1
			}

			averages[filename] = average
		}

		// Iteration order not guaranteed.
		keys := make([]string, 0, len(averages))
		for key := range averages {
			keys = append(keys, key)
		}

		fileDistances := map[string]int{}

		for i := 0; i < len(keys); i++ {
			for j := i + 1; j < len(keys); j++ {
				distance, err := fileDistance(keys[i], keys[j])
				if err != nil {
					return 0.0, err
				}
				fileDistances[keys[i]] += distance
				fileDistances[keys[j]] += distance
			}
		}

		// calculate average distance between files
		// limitations: does not weight average against file distance
		for _, key := range keys {
			fileDistance := fileDistances[key] / len(keys)
			diffDistance := averages[key] // assumption that diffDistance is not zero

			distance += diffDistance * float64(fileDistance)
		}

		return distance, nil
	}
}

func fileDistance(base, target string) (int, error) {
	path, err := filepath.Rel(base, target)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate distance between %s and %s: %w", base, target, err)
	}

	switch {
	case base == target:
		return 0, nil // same file
	case path == ".":
		return 1, nil // same directory
	default:
		// target is child of base
		return len(strings.Split(path, "/")) + 1, nil
	}
}