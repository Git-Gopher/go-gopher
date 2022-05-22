package detector

type Detector interface {
	Result() (violated, count, total int)
}
