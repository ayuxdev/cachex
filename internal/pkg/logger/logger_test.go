package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	// Test Error
	DisableWarn = false
	DisableInfo = false
	DisableDebug = false
	DisableSuccess = false
	EnableTimeStamp = true

	Info("This is an info message")
	Infof("This is an info message with format: %s", "formatted")
	Success("This is a success message")
	Successf("This is a success message with format: %s", "formatted")
	Warn("This is a warning message")
	Warnf("This is a warning message with format: %s", "formatted")
	Debug("This is a debug message")
	Debugf("This is a debug message with format: %s", "formatted")
	Errorf("This is an error message with format: %s", "formatted")
}
