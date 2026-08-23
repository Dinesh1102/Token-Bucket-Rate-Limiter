package algorithm

import (
	"testing"
	"time"

	"rate-limiter/internal/models"
)

func TestTokenBucketAllowsInitialRequest(t *testing.T) {
	now := time.Date(2026, 8, 23, 10, 0, 0, 0, time.UTC)
	tokenBucket := NewTokenBucketWithClock(func() time.Time { return now })

	decision, nextState, err := tokenBucket.Execute(models.ClientConfig{
		ClientKey: "client-1",
		Rate:      2,
		Burst:     5,
	}, models.State{})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !decision.Allowed {
		t.Fatal("expected request to be allowed")
	}
	if decision.Remaining != 4 {
		t.Fatalf("Remaining = %d, want 4", decision.Remaining)
	}
	if nextState.TokenBucket.Tokens != 4 {
		t.Fatalf("Tokens = %v, want 4", nextState.TokenBucket.Tokens)
	}
	if !nextState.TokenBucket.LastRefilled.Equal(now) {
		t.Fatalf("LastRefilled = %v, want %v", nextState.TokenBucket.LastRefilled, now)
	}
}

func TestTokenBucketDeniesWhenEmpty(t *testing.T) {
	now := time.Date(2026, 8, 23, 10, 0, 0, 0, time.UTC)
	tokenBucket := NewTokenBucketWithClock(func() time.Time { return now })

	decision, nextState, err := tokenBucket.Execute(models.ClientConfig{
		ClientKey: "client-1",
		Rate:      2,
		Burst:     5,
	}, models.State{
		TokenBucket: models.TokenBucketState{
			Tokens:       0,
			LastRefilled: now,
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if decision.Allowed {
		t.Fatal("expected request to be denied")
	}
	if decision.RetryAfter != 500*time.Millisecond {
		t.Fatalf("RetryAfter = %v, want 500ms", decision.RetryAfter)
	}
	if nextState.TokenBucket.Tokens != 0 {
		t.Fatalf("Tokens = %v, want 0", nextState.TokenBucket.Tokens)
	}
}

func TestTokenBucketRefillsBeforeAllowing(t *testing.T) {
	lastRefilled := time.Date(2026, 8, 23, 10, 0, 0, 0, time.UTC)
	now := lastRefilled.Add(2 * time.Second)
	tokenBucket := NewTokenBucketWithClock(func() time.Time { return now })

	decision, nextState, err := tokenBucket.Execute(models.ClientConfig{
		ClientKey: "client-1",
		Rate:      2,
		Burst:     5,
	}, models.State{
		TokenBucket: models.TokenBucketState{
			Tokens:       1,
			LastRefilled: lastRefilled,
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !decision.Allowed {
		t.Fatal("expected request to be allowed")
	}
	if nextState.TokenBucket.Tokens != 4 {
		t.Fatalf("Tokens = %v, want 4", nextState.TokenBucket.Tokens)
	}
	if !nextState.TokenBucket.LastRefilled.Equal(now) {
		t.Fatalf("LastRefilled = %v, want %v", nextState.TokenBucket.LastRefilled, now)
	}
}
