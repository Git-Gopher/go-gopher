package workflow

import (
	"fmt"

	"github.com/Git-Gopher/go-gopher/cache"
	"github.com/Git-Gopher/go-gopher/detector"
	"github.com/Git-Gopher/go-gopher/model/local"
	"github.com/Git-Gopher/go-gopher/violation"
)

func GithubFlowWorkflow() *Workflow {
	return &Workflow{
		Name: "Github Flow",
		WeightedCommitDetectors: []WeightedDetector{
			{Weight: 2, Detector: detector.NewCommitDetector(detector.TwoParentsCommitDetect())},
		},
		WeightedCacheDetectors: []WeightedCacheDetector{
			{Weight: 10, Detector: detector.NewCommitCacheDetector(detector.ForcePushDetect())},
		},
	}
}

type WeightedDetector struct {
	Weight   int
	Detector detector.Detector
}

type WeightedCacheDetector struct {
	Weight   int
	Detector detector.CacheDetector
}

type Workflow struct {
	Name                    string
	WeightedCommitDetectors []WeightedDetector
	WeightedCacheDetectors  []WeightedCacheDetector
}

// TODO: Use weight here.
func (w *Workflow) Analyze(model *local.GitModel) (violated int,
	count,
	total int,
	violations []violation.Violation,
	err error,
) {
	for _, wd := range w.WeightedCommitDetectors {
		if err := wd.Detector.Run(model); err != nil {
			// XXX: Change this to acceptable behavior

			return 0, 0, 0, nil, fmt.Errorf("Failed to analyze workflow: %w", err)
		}
		v, c, t, vs := wd.Detector.Result()
		violated += v
		count += c
		total += t
		violations = append(violations, vs...)
	}

	for _, wd := range w.WeightedCacheDetectors {
		caches, err := cache.ReadCaches()
		if err != nil {
			return 0, 0, 0, nil, fmt.Errorf("Failed to read caches: %w", err)
		}
		current := cache.NewCache(model)
		if err := wd.Detector.Run(current, caches); err != nil {
			// XXX: Change this to acceptable behavior

			return 0, 0, 0, nil, fmt.Errorf("Failed to analyze workflow: %w", err)
		}
		v, c, t, vs := wd.Detector.Result()
		violated += v
		count += c
		total += t
		violations = append(violations, vs...)
	}

	return violated, count, total, violations, nil
}
