package state

import (
	"sync"

	"rate-limiter/internal/models"
)

type MemoryStateManager struct {
	mu     sync.RWMutex
	states map[string]models.State
}

func NewMemoryStateManager() *MemoryStateManager {
	return &MemoryStateManager{
		states: make(map[string]models.State),
	}
}

func (m *MemoryStateManager) Load(clientKey string) (models.State, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.states[clientKey], nil
}

func (m *MemoryStateManager) Save(clientKey string, state models.State) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.states[clientKey] = state
	return nil
}
