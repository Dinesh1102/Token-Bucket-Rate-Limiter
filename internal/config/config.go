package config

import "rate-limiter/internal/models"

type ConfigManager interface {
	GetClient(clientKey string) (*models.ClientConfig, error)
}