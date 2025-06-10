package scanner

import (
	"fmt"
	"testing"

	"github.com/ayuxdev/cachex/internal/pkg/client"
	"github.com/ayuxdev/cachex/internal/pkg/config"
)

func TestRunPoisoningTest(t *testing.T) {

	labId := "0a00007f04f5a5b586fb892300d2006c"

	s := ScannerArgs{
		URL: fmt.Sprintf("https://%s.web-security-academy.net/", labId),
		RequestHeaders: map[string]string{
			"User-Agent": "Mozilla/5.0",
		},
		PayloadHeaders: map[string]string{
			"X-Forwarded-Host":   "example.com",
			"X-Forwarded-Scheme": "https",
			"X-Forwarded-Port":   "443",
			"X-Forwarded-For":    "127.0.0.1",
		},
	}

	s.PersistenceCheckerArgs = &PersistenceCheckerArgs{
		ScannerArgs:       &s,
		DoCheck:           true,
		NumRequestsToSend: 5,
		NumThreads:        5,
	}

	s.Client = config.DefaultConfig().ClientConfig.CreateNewClient()
	s.SetCacheBusterURL()

	originalResponse, err := client.FetchResponse(s.cacheBusterURL, s.RequestHeaders, s.Client)
	if err != nil {
		t.Fatalf("Error fetching original response: %v", err)
	}
	s.OriginalResponse = originalResponse

	scanResult, err := s.RunPoisoningTest()
	if err != nil {
		t.Fatalf("Error running poisoning test: %v", err)
	}

	t.Logf("Scan result: %+v", scanResult)
	if err != nil {
		t.Fatalf("Error in scan result: %v", err)
	}

	if scanResult.IsVulnerable == false {
		t.Logf("\n\n\n\n\nOriginal Response: %s", scanResult.OriginalResponse.Body)
		t.Logf("\n\n\n\n\nModified Response: %s", scanResult.ModifiedResponse.Body)
	}

	if scanResult.PersistenceCheckResult.IsPersistent == false {
		t.Logf("\n\n\n\n\nVerifier Final Response: %s", scanResult.PersistenceCheckResult.FinalResponse.Body)
		t.Logf("\n\n\n\n\nModified Response: %s", scanResult.PersistenceCheckResult.FinalResponse.Body)
	}

}
