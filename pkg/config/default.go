package config

import (
	"os"
	"path/filepath"
)

var HomeDir = os.Getenv("HOME")
var AppName = "cachex"
var DefaultCfgDir = filepath.Join(HomeDir, ".config", AppName)
var DefaultPayloadHeadersPath = filepath.Join(DefaultCfgDir, "payloads.yaml")
var DefaultScannerConfigPath = filepath.Join(DefaultCfgDir, "config.yaml")

// DefaultUserAgent is the default user agent to use for requests
var DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"

// DefaultConfig returns the default configuration for cachex
func DefaultConfig() *Config {
	return &Config{
		PayloadConfig: PayloadConfig{
			PayloadHeaders: PayloadHeaders,
		},
		ScannerConfig: ScannerConfig{
			ScanMode: "single",
			Threads:  25,
			RequestHeaders: map[string]string{
				"User-Agent": DefaultUserAgent,
				"Accept":     "*/*",
			},
			Client: ClientConfig{
				DialTimeout:      5,
				HandshakeTimeout: 5,
				ResponseTimeout:  10,
				ProxyURL:         "",
			},
			PersistenceCheckerArgs: PersistenceCheckerArgs{
				Enabled:           true,
				NumRequestsToSend: 10,
				Threads:           5,
			},
			LoggerConfig: LoggerConfig{
				LogError:     false,
				LogMode:      "pretty",
				LogTarget:    "stdout",
				Debug:        false,
				SkipTenative: true,
			},
		},
	}
}
