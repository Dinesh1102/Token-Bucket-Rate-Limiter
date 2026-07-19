package metrics

import "rate-limiter/internal/models"

type MetricsCollector interface {
	Record(decision *models.Decision)
}