package manager

import (
	"context"

	"rate-limiter/internal/algorithm"
	"rate-limiter/internal/config"
	"rate-limiter/internal/metrics"
	"rate-limiter/internal/models"
	"rate-limiter/internal/state"
)

type RateLimiterManager interface {
	Check(ctx context.Context, clientKey string) (*models.Decision, error)
}

type rateLimiterManager struct {
	configManager    config.ConfigManager
	stateManager     state.StateManager
	algorithmManager *algorithm.Manager
	metricsCollector metrics.MetricsCollector
}

func NewRateLimiterManager(
	configManager config.ConfigManager,
	stateManager state.StateManager,
	algorithmManager *algorithm.Manager,
	metricsCollector metrics.MetricsCollector,
) RateLimiterManager {
	return &rateLimiterManager{
		configManager:    configManager,
		stateManager:     stateManager,
		algorithmManager: algorithmManager,
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