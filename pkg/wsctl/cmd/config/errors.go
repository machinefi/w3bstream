package config

import "fmt"

var (
	// ErrConfigNotMatch indicates error for no config matches
	ErrConfigNotMatch = fmt.Errorf("no matching config")
	// ErrEmptyEndpoint indicates error for empty endpoint
	ErrEmptyEndpoint = fmt.Errorf("no endpoint has been set")
)
