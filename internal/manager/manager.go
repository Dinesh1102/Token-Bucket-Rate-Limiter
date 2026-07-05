package manager

import (
	"context"

	"rate-limiter/internal/models"
)

type RateLimiterManager interface {
	Check(ctx context.Context, clientKey string) (*models.Decision, error)
}