package manager

import "errors"

var (
	ErrClientConfigNotFound = errors.New("client config not found")
	ErrAlgorithmNotSet      = errors.New("algorithm not set")
)
