package workflow

import (
	"fmt"
	"log"

	"github.com/Git-Gopher/go-gopher/cache"
	"github.com/Git-Gopher/go-gopher/detector"
	"github.com/Git-Gopher/go-gopher/model/enriched"
	"github.com/Git-Gopher/go-gopher/violation"
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

func GithubFlowWorkflow() *Workflow {
	return &Workflow{
		Name: "Github Flow",
		WeightedCommitDetectors: []WeightedDetector{
			{Weight: 1, Detector: detector.NewBranchDetector(detector.StaleBranchDetect())},
			{Weight: 1, Detector: detector.NewPullRequestDetector(detector.PullRequestApprovalDetector())},
			{Weight: 1, Detector: detector.NewPullRequestDetector(detector.PullRequestIssueDetector())},
			{Weight: 1, Detector: detector.NewPullRequestDetector(detector.PullRequestReviewThreadDetector())},
			{Weight: 1, Detector: detector.NewCommitDetector(detector.DiffMatchesMessageDetect())},
			{Weight: 1, Detector: detector.NewCommitDistanceDetector(detector.DiffDistanceCalculation())},
			{Weight: 1, Detector: detector.NewBranchCompareDetector(detector.NewBranchNameConsistencyDetect())},
			{Weight: 1, Detector: detector.NewFeatureBranchDetector()},
			// DISABLED
			// {Weight: 1, Detector: detector.NewBranchCompareDetector(detector.NewFeatureBranchNameDetect())},
			// {Weight: 1, Detector: detector.NewCommitDetector(detector.TwoParentsCommitDetect())},
		},
		WeightedCacheDetectors: []WeightedCacheDetector{
			{Weight: 10, Detector: detector.NewCommitCacheDetector(detector.ForcePushDetect())},
		},
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
