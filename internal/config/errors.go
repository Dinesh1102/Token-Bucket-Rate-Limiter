package config

import "errors"

var (
	ErrClientAlreadyExists = errors.New("client config already exists")
	ErrClientNotFound      = errors.New("client config not found")
)
