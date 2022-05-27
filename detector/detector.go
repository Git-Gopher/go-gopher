package detector

import "github.com/Git-Gopher/go-gopher/model"

type Detector interface {
	Result() (violated, count, total int)
	// Eventually this will not pass the entire enriched model into the detector
	// and we will split up the commits into buckets
	Run(model *model.GitModel)
}
