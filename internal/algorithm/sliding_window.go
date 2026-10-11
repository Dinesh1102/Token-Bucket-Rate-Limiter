package algorithm

import (
	"math"
	"time"

	"rate-limiter/internal/models"
)

type SlidingWindow struct {
	now func() time.Time
}

func NewSlidingWindow() *SlidingWindow {
	return &SlidingWindow{
		now: time.Now,
	}
}

func NewSlidingWindowWithClock(now func() time.Time) *SlidingWindow {
	return &SlidingWindow{
		now: now,
	}
}

func (sw *SlidingWindow) Execute(
	config models.ClientConfig,
	state models.State,
) (*models.Decision, models.State, error) {
	if config.Rate <= 0 || config.Burst <= 0 {
		return nil, state, ErrInvalidRateLimitConfig
	}

	now := sw.now()
	windowDuration := slidingWindowDuration(config)
	window := normalizeSlidingWindow(state.SlidingWindow, now, windowDuration)
	estimatedCount := slidingWindowEstimatedCount(window, now, windowDuration)

	allowed := estimatedCount+1 <= float64(config.Burst)
	retryAfter := time.Duration(0)
	if allowed {
		window.Count++
		estimatedCount++
	} else {
		retryAfter = slidingWindowRetryAfter(window, now, windowDuration, config.Burst)
	}

	state.SlidingWindow = window

	return &models.Decision{
		Allowed:    allowed,
		Limit:      config.Burst,
		Remaining:  int(math.Max(0, math.Floor(float64(config.Burst)-estimatedCount))),
		ResetTime:  now.Add(slidingWindowResetAfter(window, now, windowDuration, config.Burst)),
		RetryAfter: retryAfter,
	}, state, nil
}

func normalizeSlidingWindow(window models.SlidingWindowState, now time.Time, windowDuration time.Duration) models.SlidingWindowState {
	if window.WindowStart.IsZero() {
		window.WindowStart = now
		return window
	}

	elapsed := now.Sub(window.WindowStart)
	if elapsed < windowDuration {
		return window
	}

	windowsElapsed := int(elapsed / windowDuration)
	if windowsElapsed == 1 {
		window.PreviousCount = window.Count
	} else {
		window.PreviousCount = 0
	}
	window.Count = 0
	window.WindowStart = window.WindowStart.Add(time.Duration(windowsElapsed) * windowDuration)

	return window
}

func slidingWindowEstimatedCount(window models.SlidingWindowState, now time.Time, windowDuration time.Duration) float64 {
	elapsed := now.Sub(window.WindowStart)
	if elapsed < 0 {
		elapsed = 0
	}
	if elapsed >= windowDuration {
		return float64(window.Count)
	}

	previousWeight := 1 - (float64(elapsed) / float64(windowDuration))
	return float64(window.Count) + (float64(window.PreviousCount) * previousWeight)
}

func slidingWindowRetryAfter(window models.SlidingWindowState, now time.Time, windowDuration time.Duration, limit int) time.Duration {
	return slidingWindowDurationUntilAtMost(window, now, windowDuration, float64(limit-1))
}

func slidingWindowResetAfter(window models.SlidingWindowState, now time.Time, windowDuration time.Duration, limit int) time.Duration {
	return slidingWindowDurationUntilAtMost(window, now, windowDuration, 0)
}

func slidingWindowDurationUntilAtMost(
	window models.SlidingWindowState,
	now time.Time,
	windowDuration time.Duration,
	targetCount float64,
) time.Duration {
	if slidingWindowEstimatedCount(window, now, windowDuration) <= targetCount {
		return 0
	}

	current := float64(window.Count)
	if window.PreviousCount <= 0 || current > targetCount {
		elapsed := now.Sub(window.WindowStart)
		if elapsed < 0 {
			elapsed = 0
		}
		return windowDuration - elapsed
	}

	previous := float64(window.PreviousCount)
	elapsed := now.Sub(window.WindowStart)
	if elapsed < 0 {
		elapsed = 0
	}

	requiredElapsed := windowDurationSeconds(windowDuration) * (1 - ((targetCount - current) / previous))
	waitSeconds := requiredElapsed - elapsed.Seconds()
	if waitSeconds <= 0 {
		return 0
	}

	return secondsToDuration(waitSeconds)
}

func slidingWindowDuration(config models.ClientConfig) time.Duration {
	return secondsToDuration(float64(config.Burst) / config.Rate)
}

func windowDurationSeconds(duration time.Duration) float64 {
	return float64(duration) / float64(time.Second)
}
