package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Test with a valid config file
	err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if Cfg == nil {
		t.Fatal("LoadConfig() did not load config")
	}
	t.Logf("LoadConfig() loaded config: %+v", Cfg)
}
