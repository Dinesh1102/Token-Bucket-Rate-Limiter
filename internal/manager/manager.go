package manager

import (
	"context"

	"rate-limiter/internal/models"
)

type RateLimiterManager interface {
	Check(ctx context.Context, clientKey string) (*models.Decision, error)
}

type rateLimiterManager struct {
	configManager    ConfigManager
	stateManager     StateManager
	algorithmFactory AlgorithmFactory
	metricsCollector MetricsCollector
}

func NewRateLimiterManager(
	configManager ConfigManager,
	stateManager StateManager,
	algorithmFactory AlgorithmFactory,
	metricsCollector MetricsCollector,
) RateLimiterManager {
	return &rateLimiterManager{
		configManager:    configManager,
		stateManager:     stateManager,
		algorithmFactory: algorithmFactory,
		metricsCollector: metricsCollector,
	}
}

func (m *rateLimiterManager) Check(
	ctx context.Context,
	clientKey string,
) (*models.Decision, error) {

	// Implementation will come next.
	panic("not implemented")
}