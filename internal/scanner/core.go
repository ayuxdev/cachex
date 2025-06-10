// Description: This file contains the core logic for the cache poisoning scanner.
package scanner

import (
	"fmt"

	"github.com/ayuxdev/cachex/internal/pkg/client"
	"github.com/ayuxdev/cachex/internal/pkg/logger"
)

// RunPoisoningTest tests the target URL for cache poisoning by:
// 1. Taking the original response without payload headers & custom payload headers to inject.
// 2. Sending a modified request with the provided payload headers to potentially poison the cache.
// 3. Comparing the modified response with the original response to detect changes, indicating response manipulation vulnerability.
// 4. Optionally checking if the changes persist across multiple requests. If both response manipulation and persistence are detected, cache poisoning is confirmed.
// It returns a ScannerOutput struct containing the results of the scan.
func (s *ScannerArgs) RunPoisoningTest() (*ScannerOutput, error) {
	s.SetCacheBusterURL() // Set the cache buster URL to ensure a fresh response

	// Merge request headers with payload headers and send the request.
	modifiedResponse, err := client.FetchResponse(s.cacheBusterURL, MergeMaps(s.RequestHeaders, s.PayloadHeaders), s.Client)
	logger.Debugf("Received modified response: %v", modifiedResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch modified response: %v", err)
	}

	// Detect changes between the original and modified responses.
	changeTypeDetected, err := DetectResponseChanges(*s.OriginalResponse, *modifiedResponse)
	if err != nil {
		return nil, fmt.Errorf("error detecting response changes: %v", err)
	}

	// Prepare output structure.
	scanResult := &ScannerOutput{
		URL:                   s.URL,
		IsResponseManipulable: changeTypeDetected != NoChange, // true if any change in response was detected after adding payload headers
		ManipulationType:      changeTypeDetected,
		RequestHeaders:        s.RequestHeaders,
		PayloadHeaders:        s.PayloadHeaders,
		OriginalResponse:      s.OriginalResponse,
		ModifiedResponse:      modifiedResponse,
	}

	// If persistence checking is enabled & response is manipulable, check if the change persists or cacheable.
	if s.PersistenceCheckerArgs.DoCheck && scanResult.IsResponseManipulable {
		changeTypeToCheck := changeTypeDetected
		persistenceCheckResult := s.PersistenceCheckerArgs.CheckPersistence(modifiedResponse, changeTypeToCheck)
		if persistenceCheckResult.Err != nil {
			// only stop the execution if the error is not a poisoning error, otherwise continue with the scan
			if _, ok := persistenceCheckResult.Err.(*PoisoningError); !ok {
				return nil, fmt.Errorf("error checking response persistence: %v", persistenceCheckResult.Err)
			}
		}

		scanResult.IsVulnerable = persistenceCheckResult.IsPersistent // vulnerable to cache poisoning if the change is persistent
		scanResult.PersistenceCheckResult = persistenceCheckResult
	}

	// TODO: Am not implementing the below code because logging with all headers might help in false positive detection
	// if scan mode is single headers we need not to log the result if more than one header is used because
	// it was sent only to save time and not to be logged
	// if s.ScanMode == SingleHeaderScanMode && len(s.PayloadHeaders) > 1 {
	// }

	if err := scanResult.Log(s.OutputFile, s.LogMode, s.LogTarget, s.SkipTenative); err != nil {
		return nil, fmt.Errorf("error logging scan result: %v", err)
	}

	return scanResult, nil
}
