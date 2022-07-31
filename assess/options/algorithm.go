package options

import log "github.com/sirupsen/logrus"

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

// ThresholdGradingAlgorithm is a configurable algorithm for calculating grade.
func ThresholdGradingAlgorithm(thresholdA, thresholdB, thresholdC int) GradingAlgorithm {
	if thresholdA < 0 || thresholdA > 100 {
		log.Fatalln("thresholdA must be between 0 and 100")
	}

	if thresholdA < thresholdB {
		log.Fatalln("thresholdA must be greater than thresholdB")
	}

	if thresholdB < 0 || thresholdB > 100 {
		log.Fatalln("thresholdB must be between 0 and 100")
	}

	if thresholdB < thresholdC {
		log.Fatalln("thresholdB must be greater than thresholdC")
	}

	if thresholdC < 0 || thresholdC > 100 {
		log.Fatalln("thresholdC must be between 0 and 100")
	}

	return func(violations int, total int) int {
		percentage := violations * 100 / total

		switch {
		case percentage < 100-thresholdA:
			return 3
		case percentage < 100-thresholdB:
			return 2
		case percentage < 100-thresholdC:
			return 1
		default:
			return 0
		}
	}
}

// GetGradingAlgorithm returns the grading algorithm.
func GetGradingAlgorithm(name string, t *ThresholdSettings) GradingAlgorithm {
	algorithMap := map[string]GradingAlgorithm{
		"basic-algorithm": BasicGradingAlgorithm,
	}

	if t != nil {
		threshold := *t

		algorithMap["threshold-algorithm"] = ThresholdGradingAlgorithm(
			threshold.ThresholdA,
			threshold.ThresholdB,
			threshold.ThresholdC,
		)
	}

	if gradingAlgorithm, ok := algorithMap[name]; ok {
		return gradingAlgorithm
	}

	log.Error("Invalid grading algorithm:", name)

	return nil
}
