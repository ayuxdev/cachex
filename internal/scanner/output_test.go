package scanner

import (
	"testing"

	"github.com/ayuxdev/cachex/internal/pkg/client"
)

func TestWriteOutputToFile(t *testing.T) {
	// Sample scan result
	scanResult := ScannerOutput{
		URL:                   "https://example.com",
		IsVulnerable:          true,
		IsResponseManipulable: true,
		ManipulationType:      ChangedLocationHeader,
		RequestHeaders:        map[string]string{"User-Agent": "Scanner"},
		PayloadHeaders:        map[string]string{"X-Forwarded-Host": "evil.com"},
		OriginalResponse: &client.Response{
			StatusCode: 200,
			Headers:    map[string][]string{"Content-Type": {"text/html"}},
			Body:       "<html>Original</html>",
			Location:   "",
		},
		ModifiedResponse: &client.Response{
			StatusCode: 302,
			Headers:    map[string][]string{"Location": {"https://evil.com"}},
			Body:       "<html>Redirected</html>",
			Location:   "https://evil.com",
		},
		PersistenceCheckResult: &PersistenceCheckResult{
			IsPersistent: true,
			PoCLink:      "https://example.com/poc",
			FinalResponse: &client.Response{
				StatusCode: 302,
				Headers:    map[string][]string{"Location": {"https://evil.com"}},
				Body:       "<html>Still Redirected</html>",
				Location:   "https://evil.com",
			},
			Err: nil,
		},
	}

	// Output file path
	outputFile := "scan_output.json"

	// marshall scan result to json
	jsonData, err := MarshalScannerOutput(scanResult, outputFile)
	if err != nil {
		t.Fatalf("failed to marshal scan result: %v", err)
	}
	// Write output to file
	err = ExportJSONToFile(jsonData, outputFile)
	if err != nil {
		t.Fatalf("failed to export JSON to file: %v", err)
	}
}
