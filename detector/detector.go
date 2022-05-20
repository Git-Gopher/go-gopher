package detector

type Detector interface {
	Result() (violated int, count int, total int)
}
