package cachex

import (
	"fmt"
)

// Validate checks the scanner config for conflicts
func (s *Scanner) Validate() error {
	// if persistence check is disabled and logging tentative results is disabled, there will be no output
	if !s.ScannerConfig.PersistenceCheckerArgs.Enabled && s.ScannerConfig.LoggerConfig.SkipTenative {
		return fmt.Errorf("no output: persistence check and tentative logging are both disabled")
	}
	if s.ScannerConfig.Threads <= 0 {
		return fmt.Errorf("invalid number of threads: %d", s.ScannerConfig.Threads)
	}
	return nil
}
