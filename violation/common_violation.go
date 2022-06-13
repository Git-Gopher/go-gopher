package violation

func NewCommonViolation(message string) Violation {
	common := &CommonViolation{display: nil, message: message}
	common.display = &display{common}
	return common
}

// Example violation.
type CommonViolation struct {
	*display
	message string
}

// Name returns the name of the Violation.
func (*CommonViolation) Name() string {
	return "CommonViolation"
}

// Message implements Violation.
func (cv *CommonViolation) Message() string {
	return cv.message
}

// FileLocation implements Violation.
func (*CommonViolation) FileLocation() (string, error) {
	return "", ErrViolationMethodNotExist
}

// LineLocation implements Violation.
func (*CommonViolation) LineLocation() (int, error) {
	return 0, ErrViolationMethodNotExist
}

// Suggestion implements Violation.
func (*CommonViolation) Suggestion() (string, error) {
	return "", ErrViolationMethodNotExist
}
