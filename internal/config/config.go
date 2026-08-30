package config

import "rate-limiter/internal/models"

type ConfigManager interface {
	CreateClient(config models.ClientConfig) error
	GetClient(clientKey string) (*models.ClientConfig, error)
	UpdateClient(config models.ClientConfig) error
	DeleteClient(clientKey string) error
}
