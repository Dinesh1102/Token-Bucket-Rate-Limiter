package algorithm

import (
	"testing"
	"time"

	"rate-limiter/internal/models"
)

func TestSlidingWindowAllowsUntilBurstLimit(t *testing.T) {
	now := time.Date(2026, 10, 11, 10, 0, 0, 0, time.UTC)
	slidingWindow := NewSlidingWindowWithClock(func() time.Time { return now })
	state := models.State{}
	config := models.ClientConfig{
		ClientKey: "client-1",
		Rate:      1,
		Burst:     2,
	}

	first, nextState, err := slidingWindow.Execute(config, state)
	if err != nil {
		t.Fatalf("Execute() first error = %v", err)
	}
	if !first.Allowed || first.Remaining != 1 {
		t.Fatalf("first decision = %+v, want allowed with 1 remaining", first)
	}

	second, nextState, err := slidingWindow.Execute(config, nextState)
	if err != nil {
		t.Fatalf("Execute() second error = %v", err)
	}
	if !second.Allowed || second.Remaining != 0 {
		t.Fatalf("second decision = %+v, want allowed with 0 remaining", second)
	}

	third, _, err := slidingWindow.Execute(config, nextState)
	if err != nil {
		t.Fatalf("Execute() third error = %v", err)
	}
	if third.Allowed {
		t.Fatal("expected third request to be denied")
	}
	if third.RetryAfter != 2*time.Second {
		t.Fatalf("RetryAfter = %v, want 2s", third.RetryAfter)
	}
}

func TestSlidingWindowCarriesPreviousWindowWeight(t *testing.T) {
	base := time.Date(2026, 10, 11, 10, 0, 0, 0, time.UTC)
	now := base.Add(4 * time.Second)
	slidingWindow := NewSlidingWindowWithClock(func() time.Time { return now })
	config := models.ClientConfig{
		ClientKey: "client-1",
		Rate:      1,
		Burst:     4,
	}
	state := models.State{
		SlidingWindow: models.SlidingWindowState{
			WindowStart: base,
			Count:       4,
		},
	}

	decision, nextState, err := slidingWindow.Execute(config, state)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if decision.Allowed {
		t.Fatal("expected request to be denied while previous window still has weight")
	}
	if nextState.SlidingWindow.PreviousCount != 4 || nextState.SlidingWindow.Count != 0 {
		t.Fatalf("next window state = %+v, want previous count 4 current count 0", nextState.SlidingWindow)
	}
	if decision.RetryAfter != 1*time.Second {
		t.Fatalf("RetryAfter = %v, want 1s", decision.RetryAfter)
	}
}

func TestSlidingWindowAllowsAfterPreviousWindowExpires(t *testing.T) {
	base := time.Date(2026, 10, 11, 10, 0, 0, 0, time.UTC)
	now := base.Add(9 * time.Second)
	slidingWindow := NewSlidingWindowWithClock(func() time.Time { return now })
	config := models.ClientConfig{
		ClientKey: "client-1",
		Rate:      1,
		Burst:     4,
	}
	state := models.State{
		SlidingWindow: models.SlidingWindowState{
			WindowStart: base,
			Count:       4,
		},
	}

	decision, nextState, err := slidingWindow.Execute(config, state)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !decision.Allowed {
		t.Fatal("expected request to be allowed after previous window expires")
	}
	if nextState.SlidingWindow.PreviousCount != 0 || nextState.SlidingWindow.Count != 1 {
		t.Fatalf("next window state = %+v, want previous count 0 current count 1", nextState.SlidingWindow)
	}
}

func TestSlidingWindowRejectsInvalidConfig(t *testing.T) {
	slidingWindow := NewSlidingWindowWithClock(func() time.Time {
		return time.Date(2026, 10, 11, 10, 0, 0, 0, time.UTC)
	})

	_, _, err := slidingWindow.Execute(models.ClientConfig{
		ClientKey: "client-1",
		Rate:      0,
		Burst:     4,
	}, models.State{})
	if err != ErrInvalidRateLimitConfig {
		t.Fatalf("Execute() error = %v, want ErrInvalidRateLimitConfig", err)
	}
}
