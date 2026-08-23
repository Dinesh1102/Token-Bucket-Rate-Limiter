package algorithm

import (
	"math"
	"time"

	"rate-limiter/internal/models"
)

type TokenBucket struct {
	now func() time.Time
}

func NewTokenBucket() *TokenBucket {
	return &TokenBucket{
		now: time.Now,
	}
}

func NewTokenBucketWithClock(now func() time.Time) *TokenBucket {
	return &TokenBucket{
		now: now,
	}
}

func (tb *TokenBucket) Execute(
	config models.ClientConfig,
	state models.State,
) (*models.Decision, models.State, error) {
	if config.Rate <= 0 || config.Burst <= 0 {
		return nil, state, ErrInvalidRateLimitConfig
	}

	now := tb.now()
	bucket := state.TokenBucket
	if bucket.LastRefilled.IsZero() {
		bucket.Tokens = float64(config.Burst)
		bucket.LastRefilled = now
	} else {
		elapsed := now.Sub(bucket.LastRefilled).Seconds()
		if elapsed > 0 {
			bucket.Tokens = math.Min(float64(config.Burst), bucket.Tokens+(elapsed*config.Rate))
			bucket.LastRefilled = now
		}
	}

	allowed := bucket.Tokens >= 1
	retryAfter := time.Duration(0)
	if allowed {
		bucket.Tokens--
	} else {
		retryAfter = secondsToDuration((1 - bucket.Tokens) / config.Rate)
	}

	state.TokenBucket = bucket

	return &models.Decision{
		Allowed:    allowed,
		Limit:      config.Burst,
		Remaining:  int(math.Floor(bucket.Tokens)),
		ResetTime:  resetTime(now, bucket.Tokens, config),
		RetryAfter: retryAfter,
	}, state, nil
}

func resetTime(now time.Time, tokens float64, config models.ClientConfig) time.Time {
	missingTokens := float64(config.Burst) - tokens
	if missingTokens <= 0 {
		return now
	}

	return now.Add(secondsToDuration(missingTokens / config.Rate))
}

func secondsToDuration(seconds float64) time.Duration {
	return time.Duration(math.Ceil(seconds * float64(time.Second)))
}
