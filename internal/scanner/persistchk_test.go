package scanner

import (
	"fmt"
	"testing"

	"github.com/ayuxdev/cachex/internal/pkg/client"
	"github.com/ayuxdev/cachex/internal/pkg/config"
)

func TestVerifyResponsePersistence(t *testing.T) {
	labId := "0a00007f04f5a5b586fb892300d2006c"

	s := ScannerArgs{
		URL: fmt.Sprintf("https://%s.web-security-academy.net/", labId),
		RequestHeaders: map[string]string{
			"User-Agent": "Mozilla/5.0",
		},
		PayloadHeaders: map[string]string{
			"X-Forwarded-Host": "example.com",
		},
	}

	s.PersistenceCheckerArgs = &PersistenceCheckerArgs{
		ScannerArgs:       &s,
		DoCheck:           true,
		NumRequestsToSend: 5,
		NumThreads:        5,
	}

	s.Client = config.DefaultConfig().ClientConfig.CreateNewClient()

	modifiedResponse, err := client.FetchResponse(s.URL, MergeMaps(s.PayloadHeaders, s.RequestHeaders), s.Client)
	if err != nil {
		t.Fatalf("error fetching response: %v", err)
	}

	changeTypeToVerify := ChangedBody

	PersistenceCheckResult := s.PersistenceCheckerArgs.CheckPersistence(modifiedResponse, changeTypeToVerify)
	if PersistenceCheckResult.Err != nil {
		t.Fatalf("error occurred while verifying response persistence: %v", PersistenceCheckResult.Err)
	}

	t.Log(PersistenceCheckResult)
}
