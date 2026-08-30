package metrics

import "rate-limiter/internal/models"

type NoopCollector struct{}

func NewNoopCollector() *NoopCollector {
	return &NoopCollector{}
}

func (c *NoopCollector) Record(decision *models.Decision) {}
