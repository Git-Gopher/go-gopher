package workflow

import "github.com/Git-Gopher/go-gopher/detector"

type Workflow struct {
	name      string
	detectors []detector.Detector
}

// XXX: Do foldr instead
func (w *Workflow) Analyze() (violated int, count int, total int) {
	for _, d := range w.detectors {
		v, c, t := d.Result()
		violated += v
		count += c
		total += t
	}

	return
}
