package detector

import "github.com/Git-Gopher/go-gopher/model"

type Detector interface {
	Result() (violated, count, total int)
	Run(model *model.GitModel) error
}
