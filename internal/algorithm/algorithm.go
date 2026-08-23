package algorithm

import "rate-limiter/internal/models"

type Algorithm interface {
	Execute(
		config models.ClientConfig,
		state models.State,
	) (*models.Decision, models.State, error)
}

type Manager struct {
	tokenBucket   Algorithm
	slidingWindow Algorithm
}

func NewManager(
	tokenBucket Algorithm,
	slidingWindow Algorithm,
) *Manager {
	return &Manager{
		tokenBucket:   tokenBucket,
		slidingWindow: slidingWindow,
	}
}

func (m *Manager) Get(name models.AlgorithmName) (Algorithm, error) {

	switch name {

	case models.AlgorithmTokenBucket:
		return m.tokenBucket, nil

	case models.AlgorithmSlidingWindow:
		return m.slidingWindow, nil

	default:
		return nil, ErrUnknownAlgorithm
	}
}
