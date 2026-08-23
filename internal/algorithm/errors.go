package algorithm

import "errors"

var ErrUnknownAlgorithm = errors.New("unknown algorithm")
var ErrInvalidRateLimitConfig = errors.New("invalid rate limit config")
