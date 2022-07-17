package assess

// All grades are out of 3 points.
type GradingAlgorithm func(violations int, total int) int

// BasicGradingAlgorithm is a basic algorithm for calculating grade.
// All grades are out of 3 points.
func BasicGradingAlgorithm(violations int, total int) int {
	percentage := violations * 100 / total

	switch {
	case percentage < 20:
		return 3 // 100% > x > 80%
	case percentage < 40:
		return 2 // 80% > x > 60%
	case percentage < 60:
		return 1 // 60% > x > 40%
	default:
		return 0 // 40% > x > 0%
	}
}
