package detector

import (
	"github.com/Git-Gopher/go-gopher/model/local"
)

type Detector interface {
	Result() (violated, count, total int)
	// TODO: We should change this to the enriched model
	Run(model *local.GitModel) error
}
