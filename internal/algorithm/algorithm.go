package algorithm

import "rate-limiter/internal/models"

type Algorithm interface {
	Execute(
		config models.ClientConfig,
		state models.State,
	) (*models.Decision, models.State, error)
}

type Manager struct {
	tokenBucket  Algorithm
	slidingWindow Algorithm
}

func NewManager(
	tokenBucket Algorithm,
	slidingWindow Algorithm,
) *Manager {
	return &Manager{
		tokenBucket:  tokenBucket,
		slidingWindow: slidingWindow,
	}
}

func (m *Manager) Get(name string) (Algorithm, error) {

	switch name {

	case "TOKEN_BUCKET":
		return m.tokenBucket, nil

	case "SLIDING_WINDOW":
		return m.slidingWindow, nil

	default:
		return nil, ErrUnknownAlgorithm
	}
}