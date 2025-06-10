// Description: This file contains the implementation of the main scanner which is responsible for scanning the target URL for cache poisoning vulnerabilities.
package scanner

import (
	"fmt"

	"github.com/ayuxdev/cachex/internal/pkg/client"
	log "github.com/ayuxdev/cachex/internal/pkg/logger"
)

// RunBatchScan runs the cache poisoning scan on a batch of URLs concurrently and returns results + all encountered errors.
func (s *ScannerArgs) RunBatchScan(URLs []string, Threads int) ([]ScannerOutput, []error) {
	type result struct {
		output []ScannerOutput
		err    error
	}

	resultsChan := make(chan result, len(URLs))
	sem := make(chan struct{}, Threads)

	for _, url := range URLs {
		sem <- struct{}{} // acquire slot

		go func(target string) {
			defer func() { <-sem }() // release slot

			// Clone scanner args
			args := *s
			args.URL = target

			// Clone persistence checker properly
			if s.PersistenceCheckerArgs != nil {
				clone := *s.PersistenceCheckerArgs
				clone.ScannerArgs = &args
				args.PersistenceCheckerArgs = &clone
			}

			scanResult, err := args.Run()

			if err != nil && args.LogError {
				log.Errorf("failed to scan endpoint %s: %v", target, err)
			}
			resultsChan <- result{output: scanResult, err: err}
		}(url)
	}

	var finalResults []ScannerOutput
	var allErrors []error

	for range URLs {
		res := <-resultsChan
		if res.err != nil {
			allErrors = append(allErrors, res.err)
		}
		if res.output != nil {
			finalResults = append(finalResults, res.output...)
		}
	}

	return finalResults, allErrors
}

// Run is the main function that performs the cache poisoning scan.
// It first fetches the original response without payload headers and tests each payload header separately or all together based on the scan mode.
// The original response is then passed to `RunPoisoningTest` with payload headers to detect changes which is the core function.
// It returns a slice of ScannerOutput structs containing the results of the scan.
func (s *ScannerArgs) Run() ([]ScannerOutput, error) {
	// Skipping the retrieval of `originalResponse` without `CacheBusterURL` to prevent unintended variations in the response.
	s.SetCacheBusterURL()
	// Fetch original response without payload headers
	originalResponse, err := client.FetchResponse(s.cacheBusterURL, s.RequestHeaders, s.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch response: %v", err)
	}
	log.Debugf("Received original response: %v", originalResponse)
	s.OriginalResponse = originalResponse

	if s.ScanMode == SingleHeaderScanMode {
		// first scan with all payload headers eitherway to check if the response is manipulable
		// this will save time if there is no change detected
		originalDoCheck := s.PersistenceCheckerArgs.DoCheck                   // store the original value
		s.PersistenceCheckerArgs.DoCheck = false                              // temporarily override persistence checker to false since we are only checking for manipulability
		defer func() { s.PersistenceCheckerArgs.DoCheck = originalDoCheck }() // reset persistence checker to the original arg value after the scan
		multiHeaderScanResult, err := s.RunPoisoningTest()
		if err != nil {
			return nil, fmt.Errorf("error running multi-header poisoning test: %v", err)
		}

		if !multiHeaderScanResult.IsResponseManipulable {
			return []ScannerOutput{*multiHeaderScanResult}, nil
		}

		s.PersistenceCheckerArgs.DoCheck = originalDoCheck // reset persistence checker to the original arg value

		var results []ScannerOutput // Slice to hold scan results for each header

		// Test each payload header separately
		for header, value := range s.PayloadHeaders {

			s.PayloadHeaders = map[string]string{header: value} // Reset payload headers to test one header at a time

			scanResult, err := s.RunPoisoningTest()
			if err != nil {
				return nil, fmt.Errorf("error running single-header poisoning test for payload header %s: %v", header, err)
			}

			// append the scan result to the results slice
			results = append(results, *scanResult)
		}

		return results, nil
	}

	// Multi-header scan (all payload headers at once)
	scanResult, err := s.RunPoisoningTest()
	if err != nil {
		return nil, fmt.Errorf("error running multi-header poisoning test: %v", err)
	}

	return []ScannerOutput{*scanResult}, nil
}
