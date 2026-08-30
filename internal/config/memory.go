package config

import (
	"sync"

	"rate-limiter/internal/models"
)

type MemoryConfigManager struct {
	mu      sync.RWMutex
	clients map[string]models.ClientConfig
}

func NewMemoryConfigManager() *MemoryConfigManager {
	return &MemoryConfigManager{
		clients: make(map[string]models.ClientConfig),
	}
}

func (m *MemoryConfigManager) CreateClient(config models.ClientConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.clients[config.ClientKey]; exists {
		return ErrClientAlreadyExists
	}

	m.clients[config.ClientKey] = config
	return nil
}

func (m *MemoryConfigManager) GetClient(clientKey string) (*models.ClientConfig, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	config, exists := m.clients[clientKey]
	if !exists {
		return nil, ErrClientNotFound
	}

	return &config, nil
}

func (m *MemoryConfigManager) UpdateClient(config models.ClientConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.clients[config.ClientKey]; !exists {
		return ErrClientNotFound
	}

	m.clients[config.ClientKey] = config
	return nil
}

func (m *MemoryConfigManager) DeleteClient(clientKey string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.clients[clientKey]; !exists {
		return ErrClientNotFound
	}

	delete(m.clients, clientKey)
	return nil
}
