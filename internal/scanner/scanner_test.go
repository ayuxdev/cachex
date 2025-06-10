package scanner

import (
	"fmt"
	"testing"

	"github.com/ayuxdev/cachex/internal/pkg/config"
)

func TestScanner(t *testing.T) {
	labId := "0a1c002d03f81e1982b8bb6b00510041"

	// Properly initialize the ScannerArgs
	s := &ScannerArgs{
		URL: fmt.Sprintf("https://%s.web-security-academy.net/", labId),
		RequestHeaders: map[string]string{
			"User-Agent": "Mozilla/5.0",
		},
		PayloadHeaders: map[string]string{
			"X-Forwarded-Host":   "example.com",
			"X-Forwarded-Scheme": "https",
			"X-Forwarded-For":    "127.0.0.1",
		},
		ScanMode: SingleHeaderScanMode,
	}

	// Initialize PersistenceCheckerArgs
	s.PersistenceCheckerArgs = &PersistenceCheckerArgs{
		ScannerArgs:       s, // Attach ScannerArgs
		DoCheck:           true,
		NumRequestsToSend: 5,
		NumThreads:        5,
	}

	// Ensure the client is properly initialized
	s.Client = config.DefaultConfig().ClientConfig.CreateNewClient()

	// Run the scan
	scanResult, err := s.Run()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Log results
	for _, result := range scanResult {
		t.Logf("%s: %+v", result.PayloadHeaders, result)
	}
}
