package workflow

import (
	"github.com/Git-Gopher/go-gopher/detector"
	"github.com/Git-Gopher/go-gopher/model"
)

// Export some configurations for known detectors, not sure if these should be their own file or if they are configuration objects. Seems a bit strange to just have a file with this in it
var (
	GithubFlowWorkflow = Workflow{
		Name: "Github Flow",
		WeightedDetectors: []WeightedDetector{
			{Weight: 2, Detector: detector.NewCommitDetector(detector.TwoParentsCommitDetect())},
		},
	}
)

type WeightedDetector struct {
	Weight   int
	Detector detector.Detector
}
type Workflow struct {
	Name              string
	WeightedDetectors []WeightedDetector
}

// XXX: Do foldr instead
func (w *Workflow) Analyze(model *model.GitModel) (violated int, count int, total int) {
	for _, wd := range w.WeightedDetectors {
		v, c, t := wd.Detector.Result()
		violated += v
		count += c
		total += t
	}

	return
}
