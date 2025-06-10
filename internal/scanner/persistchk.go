package scanner

import (
	"fmt"
	"sync"

	"github.com/ayuxdev/cachex/internal/pkg/client"
)

// PoisoningError aggregates multiple poisoning errors
type PoisoningError struct {
	Errors []error
}

// Error implements the error interface for PoisoningError
func (pe *PoisoningError) Error() string {
	if len(pe.Errors) == 0 {
		return ""
	}
	errMsg := "Errors during poisoning:\n"
	for _, err := range pe.Errors {
		errMsg += fmt.Sprintf(" - %v\n", err)
	}
	return errMsg
}

// CheckPersistence determines whether the modified response persists across requests.
func (p *PersistenceCheckerArgs) CheckPersistence(
	modifiedResponse *client.Response,
	changeTypeToCheck ResponseChangeType) *PersistenceCheckResult {

	// If checking is disabled, return false
	if !p.DoCheck {
		return &PersistenceCheckResult{IsPersistent: false}
	}

	// Set a new cache buster URL to avoid receiving an already cached response
	p.SetCacheBusterURL()

	// Channel to collect errors without stopping execution
	errChan := make(chan error, p.NumRequestsToSend)
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, p.NumThreads) // Limit concurrency

	// Send multiple requests concurrently
	for i := range p.NumRequestsToSend {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(attempt int) {
			defer wg.Done()
			defer func() { <-semaphore }()

			err := client.SendRequest(
				p.cacheBusterURL,
				MergeMaps(p.RequestHeaders, p.PayloadHeaders),
				p.Client,
			)
			if err != nil {
				errChan <- fmt.Errorf("attempt %d: %v", attempt+1, err)
			}
		}(i)
	}

	wg.Wait()
	close(errChan) // Close channel after all requests finish

	// Collect all errors
	var poisoningErrors []error
	for err := range errChan {
		poisoningErrors = append(poisoningErrors, err)
	}

	// Fetch response without payload headers to check persistence
	response, err := client.FetchResponse(p.cacheBusterURL, p.RequestHeaders, p.Client)
	if err != nil {
		return &PersistenceCheckResult{Err: fmt.Errorf("error while fetching response without payload headers: %v", err)}
	}

	result := &PersistenceCheckResult{} // Initialize the result struct

	// Verify if the final response matches the modified response indicating persistence.
	switch changeTypeToCheck {
	case ChangedLocationHeader:
		if response.Location == modifiedResponse.Location {
			result.IsPersistent = true
		}
	case ChangedStatusCode:
		if response.StatusCode == modifiedResponse.StatusCode {
			result.IsPersistent = true
		}
	case ChangedBody:
		if response.Body == modifiedResponse.Body {
			result.IsPersistent = true
		}
	}

	result.FinalResponse = response

	// If persistence is detected, set PoC link
	if result.IsPersistent {
		result.PoCLink = p.cacheBusterURL
	}

	// If there were poisoning errors, wrap them in a custom error
	if len(poisoningErrors) > 0 {
		result.Err = &PoisoningError{Errors: poisoningErrors}
	}

	return result
}
