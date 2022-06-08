package detector

import "github.com/Git-Gopher/go-gopher/model"

type Detector interface {
	Result() (violated, count, total int)
	// TODO: We should change this to the enriched model
	Run(model *model.GitModel) error
}
