package state

import "rate-limiter/internal/models"

type StateManager interface {
	Load(clientKey string) (models.State, error)
	Save(clientKey string, state models.State) error
}