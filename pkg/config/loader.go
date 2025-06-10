package config

import (
	"fmt"
	"os"

	log "github.com/ayuxdev/cachex/internal/pkg/logger"
	"gopkg.in/yaml.v3"
)

// Cfg is the global configuration object
var Cfg *Config = DefaultConfig()

// LoadConfig loads the configuration from the config file and sets the global Cfg object
func LoadConfig() error {
	if err := os.MkdirAll(DefaultCfgDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}
	// log.Infof("Created Configuration directory: %s", DefaultCfgDir)
	if err := ensurePayloadHeadersConfig(); err != nil {
		return err
	}
	if err := ensureScannerConfig(); err != nil {
		return err
	}
	return nil
}

func ensurePayloadHeadersConfig() error {
	path := DefaultPayloadHeadersPath

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Infof("created %s", path)
		if err := SaveDefaultPayloadHeadersConfig(); err != nil {
			return fmt.Errorf("failed to save default payload headers config: %v", err)
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read payload headers config: %v", err)
	}

	return yaml.Unmarshal(data, &Cfg.PayloadConfig)
}

func ensureScannerConfig() error {
	path := DefaultScannerConfigPath

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Infof("created %s", path)
		if err := SaveDefaultScannerConfig(); err != nil {
			return fmt.Errorf("failed to save default scanner config: %v", err)
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read scanner config: %v", err)
	}

	return yaml.Unmarshal(data, &Cfg.ScannerConfig)
}

func SaveDefaultPayloadHeadersConfig() error {
	data, err := yaml.Marshal(Cfg.PayloadConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal payload headers config: %v", err)
	}
	return os.WriteFile(DefaultPayloadHeadersPath, data, 0644)
}

func SaveDefaultScannerConfig() error {
	data, err := yaml.Marshal(Cfg.ScannerConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal scanner config: %v", err)
	}
	return os.WriteFile(DefaultScannerConfigPath, data, 0644)
}
