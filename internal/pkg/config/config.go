package config

import (
	"time"

	"github.com/ayuxdev/cachex/internal/pkg/client"
)

// Config represents application wide configuration
type Config struct {
	ClientConfig client.Config
}

// DefaultConfig returns a default configuration for the application
func DefaultConfig() Config {
	config := Config{
		ClientConfig: client.Config{
			DialTimeout:           5 * time.Second,
			HandshakeTimeout:      5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
		},
	}
	return config
}
