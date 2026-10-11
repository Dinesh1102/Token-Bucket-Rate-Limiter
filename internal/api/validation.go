package api

import (
	"fmt"

	"rate-limiter/internal/models"
)

func validateClientConfig(algorithm string, rate float64, burst int) error {
	if !isSupportedAlgorithm(models.AlgorithmName(algorithm)) {
		return fmt.Errorf("unsupported algorithm %q", algorithm)
	}
	if rate <= 0 {
		return fmt.Errorf("rate must be greater than 0")
	}
	if burst <= 0 {
		return fmt.Errorf("burst must be greater than 0")
	}

	return nil
}

func isSupportedAlgorithm(algorithm models.AlgorithmName) bool {
	switch algorithm {
	case models.AlgorithmTokenBucket, models.AlgorithmSlidingWindow:
		return true
	default:
		return false
	}
}
