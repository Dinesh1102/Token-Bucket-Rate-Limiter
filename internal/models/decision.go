package models

import (
	"strconv"
	"time"
)

type Decision struct {
	Allowed    bool
	Limit      int
	Remaining  int
	ResetTime  time.Time
	RetryAfter time.Duration
}

func (d Decision) LimitString() string {
	return strconv.Itoa(d.Limit)
}

func (d Decision) RemainingString() string {
	return strconv.Itoa(d.Remaining)
}

func (d Decision) ResetString() string {
	return strconv.FormatInt(d.ResetTime.Unix(), 10)
}

func (d Decision) RetryAfterString() string {
	return strconv.Itoa(int(d.RetryAfter.Seconds()))
}