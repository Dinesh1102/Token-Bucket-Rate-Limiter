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
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	clientConfig, err := m.configManager.GetClient(clientKey)
	if err != nil {
		return nil, err
	}
	if clientConfig == nil {
		return nil, ErrClientConfigNotFound
	}

	currentState, err := m.stateManager.Load(clientKey)
	if err != nil {
		return nil, err
	}

	rateLimitAlgorithm, err := m.algorithmManager.Get(clientConfig.Algorithm)
	if err != nil {
		return nil, err
	}
	if rateLimitAlgorithm == nil {
		return nil, ErrAlgorithmNotSet
	}

	decision, nextState, err := rateLimitAlgorithm.Execute(*clientConfig, currentState)
	if err != nil {
		return nil, err
	}

	if err := m.stateManager.Save(clientKey, nextState); err != nil {
		return nil, err
	}

	if m.metricsCollector != nil {
		m.metricsCollector.Record(decision)
	}

	return decision, nil
}
