// Description: This file defines the types and structures used for cache poisoning scanning.
// It includes the ScannerArgs struct for input parameters, ScannerOutput struct for output results,
// ScanMode enum for different scanning modes, and ResponseChangeType enum for different response modification types.
package scanner

import (
	"net/http"

	"github.com/ayuxdev/cachex/internal/pkg/client"
)

// ScannerArgs defines the parameters required for scanning a target URL for cache poisoning vulnerabilities.
type ScannerArgs struct {
	URL                    string                  // Target URL to scan
	cacheBusterURL         string                  // Internal
	ScanMode               ScanMode                // Mode of scanning (single or multi-header)
	RequestHeaders         map[string]string       // Headers to be sent with the request
	PayloadHeaders         map[string]string       // Headers to be used for cache poisoning
	PersistenceCheckerArgs *PersistenceCheckerArgs // Arguments for checking cache persistence
	OriginalResponse       *client.Response        // Original response without payload headers
	Client                 *http.Client            // HTTP client to be used for sending requests
	LoggerArgs
}

type LoggerArgs struct {
	LogError     bool      // Flag to log errors to stderr
	LogMode      LogMode   // Mode of logging (pretty or JSON)
	LogTarget    LogTarget // Log target (stdout, file, both)
	OutputFile   string    // File to write the output to (optional)
	SkipTenative bool      // Skip logging tentative vulnerabilities
}

// ScannerOutput holds the results of a cache poisoning scan.
type ScannerOutput struct {
	URL                    string                  `json:"URL"`                    // Target URL that was scanned
	IsVulnerable           bool                    `json:"IsVulnerable"`           // If vulnerable to cache poisoning
	IsResponseManipulable  bool                    `json:"IsResponseManipulable"`  // If response is manipulable
	ManipulationType       ResponseChangeType      `json:"ManipulationType"`       // Type of vulnerability detected
	RequestHeaders         map[string]string       `json:"RequestHeaders"`         // Headers sent with the request
	PayloadHeaders         map[string]string       `json:"PayloadHeaders"`         // Headers used for cache poisoning
	OriginalResponse       *client.Response        `json:"OriginalResponse"`       // Original response without payload headers
	ModifiedResponse       *client.Response        `json:"ModifiedResponse"`       // Modified response with payload headers
	PersistenceCheckResult *PersistenceCheckResult `json:"PersistenceCheckResult"` // Result of persistence check
}

// PersistenceCheckerArgs contains parameters for checking response persistence.
type PersistenceCheckerArgs struct {
	*ScannerArgs           // Embedding ScannerArgs to reuse its fields
	DoCheck           bool // Whether to check response persistence
	NumRequestsToSend int  // Number of requests to send for poisoning
	NumThreads        int  // Number of threads to use for poisoning
}

// PersistenceCheckResult represents the result of checking response manipulation persistence.
type PersistenceCheckResult struct {
	IsPersistent  bool             `json:"IsPersistent"`  // Whether the modified response persists across requests
	PoCLink       string           `json:"PoCLink"`       // Link to the URL demonstrating response persistence
	FinalResponse *client.Response `json:"FinalResponse"` // Final response without payload headers after checking persistence
	Err           error            `json:"Err,omitempty"` // Error occurred during persistence check
}

// Scan Mode
type ScanMode int

// Scan modes for cache poisoning
const (
	SingleHeaderScanMode ScanMode = iota // Scan with each payload header separately
	MultiHeaderScanMode                  // Scan with all payload headers together
)

// ResponseChangeType represents different types of response modifications.
type ResponseChangeType int

// Response modification types
const (
	ChangedLocationHeader ResponseChangeType = iota // Location header changed
	ChangedStatusCode                               // Status code changed
	ChangedBody                                     // Response body changed
	NoChange                                        // No change detected btwn the responses
)
