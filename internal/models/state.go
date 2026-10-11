package models

import "time"

type State struct {
	TokenBucket   TokenBucketState
	SlidingWindow SlidingWindowState
}

type TokenBucketState struct {
	Tokens       float64
	LastRefilled time.Time
}

type SlidingWindowState struct {
	WindowStart   time.Time
	Count         int
	PreviousCount int
}
