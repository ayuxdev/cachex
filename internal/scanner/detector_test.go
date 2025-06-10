package scanner

import (
	"testing"

	"github.com/ayuxdev/cachex/internal/pkg/client"
)

func TestDetectResponseChanges(t *testing.T) {
	originalResponse := client.Response{
		Location:   "",
		StatusCode: 200,
		Body:       "Hello, World!",
	}

	modifiedResponse := client.Response{
		Location:   "",
		StatusCode: 200,
		Body:       "Hello, Pwned?",
	}

	changeType, err := DetectResponseChanges(originalResponse, modifiedResponse)

	if err != nil {
		t.Fatalf("Error occurred while detecting response changes: %v", err)
	}

	switch changeType {
	case ChangedLocationHeader:
		t.Log("Location header changed")
	case ChangedStatusCode:
		t.Log("Status code changed")
	case ChangedBody:
		t.Log("Body changed")
	case NoChange:
		t.Log("No change detected")
	}
}
