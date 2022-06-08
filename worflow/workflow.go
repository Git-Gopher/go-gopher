package workflow

import (
	"fmt"

	"github.com/Git-Gopher/go-gopher/detector"
	"github.com/Git-Gopher/go-gopher/model"
)

func GithubFlowWorkflow() *Workflow {
	return &Workflow{
		Name: "Github Flow",
		WeightedDetectors: []WeightedDetector{
			{Weight: 2, Detector: detector.NewCommitDetector(detector.TwoParentsCommitDetect())},
			{Weight: 1, Detector: detector.NewPullRequestDetector(detector.NewPullRequestIssueDetector())},
		},
	}
}

type WeightedDetector struct {
	Weight   int
	Detector detector.Detector
}

type Workflow struct {
	Name              string
	WeightedDetectors []WeightedDetector
}

// TODO: Use weight here.
func (w *Workflow) Analyze(model *model.GitModel) (violated int, count int, total int, err error) {
	for _, wd := range w.WeightedDetectors {
		if err := wd.Detector.Run(model); err != nil {
			// XXX: Change this to acceptable behavior

			return 0, 0, 0, fmt.Errorf("Failed to analyze workflow: %w", err)
		}
		v, c, t := wd.Detector.Result()
		violated += v
		count += c
		total += t
	}

	return
}
