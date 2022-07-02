package workflow

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/Git-Gopher/go-gopher/cache"
	"github.com/Git-Gopher/go-gopher/config"
	"github.com/Git-Gopher/go-gopher/detector"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/violation"
)

var (
	DefaultCsvPath   = "summary.csv"
	DetectorRegistry = map[string]detector.Detector{
		"StaleBranchDetect": detector.NewBranchDetector(detector.StaleBranchDetect()),
	}
	CacheDetectorRegistry = map[string]detector.CacheDetector{
		"ForcePush": detector.NewBranchDetector(detector.StaleBranchDetect()),
	}
)

type Workflow struct {
	Name                    string
	WeightedCommitDetectors []WeightedDetector
	WeightedCacheDetectors  []WeightedCacheDetector
}

type WeightedDetector struct {
	Weight   int
	Detector detector.Detector
}

type WeightedCacheDetector struct {
	Weight   int
	Detector detector.CacheDetector
}

func GithubFlowWorkflow(cfg *config.Config) *Workflow {
	comd, cacd := preprocess(cfg)
	return &Workflow{
		Name:                    "Github Flow",
		WeightedCommitDetectors: comd,
		WeightedCacheDetectors:  cacd,
	}
}

// LocalDetectors are detectors that can run locally without GitHub API calls.
func LocalDetectors() []detector.Detector {
	return []detector.Detector{
		detector.NewBranchDetector(detector.StaleBranchDetect()),
		detector.NewCommitDetector(detector.DiffMatchesMessageDetect()),
		detector.NewCommitDetector(detector.ShortCommitMessageDetect()),
		detector.NewCommitDistanceDetector(detector.DiffDistanceCalculation()),
		detector.NewBranchCompareDetector(detector.NewBranchNameConsistencyDetect()),
		detector.NewCommitDetector(detector.BranchCommitDetect()), // used to check if branches are used
		detector.NewFeatureBranchDetector(),
		detector.NewBranchMatrixDetector(detector.NewCrissCrossMergeDetect()),
	}
}

// Run analysis on the git project for all the detectors defined by the workflow.
func (w *Workflow) Analyze(model *enriched.EnrichedModel, current *cache.Cache, caches []*cache.Cache) (violated int,
	count,
	total int,
	violations []violation.Violation,
	err error,
) {
	v, c, t, vs, err := w.RunWeightedDetectors(model)
	if err != nil {
		return 0, 0, 0, nil, fmt.Errorf("Failed to analyze workflow: %w", err)
	}
	add(&violated, &count, &total, &violations, v, c, t, &vs, 1)

	// Only run when we have a cache
	if current != nil && caches != nil {
		v, c, t, vs, err = w.RunCacheDetectors(current, caches)
		if err != nil {
			return 0, 0, 0, nil, fmt.Errorf("Failed to analyze workflow: %w", err)
		}
		add(&violated, &count, &total, &violations, v, c, t, &vs, 1)
	} else {
		log.Println("No cache loaded, skipping cache detectors")
	}

	return violated, count, total, violations, nil
}

func (w *Workflow) RunWeightedDetectors(model *enriched.EnrichedModel) (
	violated,
	count,
	total int,
	violations []violation.Violation,
	err error,
) {
	for _, wd := range w.WeightedCommitDetectors {
		if err := wd.Detector.Run(model); err != nil {
			return 0, 0, 0, nil, fmt.Errorf("Failed to run weighted detectors: %w", err)
		}

		v, c, t, vs := wd.Detector.Result()
		add(&violated, &count, &total, &violations, v, c, t, &vs, wd.Weight)
	}

	return
}

// All cache detectors share the same current and cache, treated as readonly.
func (w *Workflow) RunCacheDetectors(current *cache.Cache, caches []*cache.Cache) (
	int,
	int,
	int,
	[]violation.Violation,
	error,
) {
	violated, count, total := 0, 0, 0
	violations := []violation.Violation{}
	for _, wd := range w.WeightedCacheDetectors {
		if err := wd.Detector.Run(current, caches); err != nil {
			return 0, 0, 0, nil, fmt.Errorf("Failed to analyze caches: %w", err)
		}

		v, c, t, vs := wd.Detector.Result()
		add(&violated, &count, &total, &violations, v, c, t, &vs, wd.Weight)
	}

	// No violations means we can reset cache to current, otherwise append to cache
	var nc []*cache.Cache
	if len(violations) == 0 {
		nc = []*cache.Cache{current}
	} else {
		nc = append(nc, caches...)
		nc = append(nc, current)
	}
	if err := cache.WriteCaches(nc); err != nil {
		return 0, 0, 0, nil, fmt.Errorf("Failed to write cache: %w", err)
	}

	return violated, count, total, violations, nil
}

// Add weighted result a detector to the shared result.
func add(
	violated, count, total *int, violations *[]violation.Violation,
	v, c, t int, vs *[]violation.Violation,
	weight int,
) {
	*violated += v * weight
	*count += c * weight
	*total += t * weight
	*violations = append(*violations, *vs...)
}

// Summarize the results of the analysis into a csv file.
func (wk *Workflow) Csv(path string) error {
	fh, err := os.Create(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("Failed to export workflow to csv file: %w", err)
	}

	// nolint
	defer fh.Close()
	// var rows [][]string
	// Construct header
	header := []string{"Repository", "URL"}
	// var body []string
	for _, wd := range wk.WeightedCommitDetectors {
		t := reflect.TypeOf(wd.Detector)
		header = append(header, t.Elem().Name())
		// violated, _, _, _ := wd.Detector.Result()
		// body = append(violated.String())
	}

	for _, wcd := range wk.WeightedCacheDetectors {
		t := reflect.TypeOf(wcd.Detector)
		header = append(header, t.Elem().Name())
	}

	w := csv.NewWriter(fh)
	err = w.Write(header)

	fmt.Printf("header: %v\n", header)

	if err != nil {
		return fmt.Errorf("Failed to write header: %w", err)
	}

	return nil
}

// Enable or disable detectors based on config.
func preprocess(cfg *config.Config) ([]WeightedDetector, []WeightedCacheDetector) {
	return nil, nil
}
