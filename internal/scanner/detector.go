// Description: This file contains functions used to detect response changes between the original and modified responses.
package scanner

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/ayuxdev/cachex/internal/pkg/client"
	"github.com/ayuxdev/cachex/internal/pkg/logger"
)

// DetectResponseChanges analyzes the differences between the original and modified responses.
func DetectResponseChanges(
	originalResponse client.Response,
	modifiedResponse client.Response) (ResponseChangeType, error) {
	// check if the responses are empty i.e. default values
	if isEmptyResponse(originalResponse) || isEmptyResponse(modifiedResponse) {
		return 0, fmt.Errorf("original or modified response is empty")
	}

	// Normalize the responses for comparison
	originalResponse = normalizeResponse(originalResponse)
	logger.Debugf("Normalized original response: %v", originalResponse)
	modifiedResponse = normalizeResponse(modifiedResponse)
	logger.Debugf("Normalized modified response: %v", modifiedResponse)

	var changeType ResponseChangeType
	// compare the responses
	switch {
	case modifiedResponse.Location != originalResponse.Location:
		changeType = ChangedLocationHeader
	case modifiedResponse.StatusCode != originalResponse.StatusCode && modifiedResponse.StatusCode != 429: // 429 is a special case for rate limiting, we ignore it as false positive
		changeType = ChangedStatusCode
	case modifiedResponse.Body != originalResponse.Body:
		changeType = ChangedBody
	default:
		changeType = NoChange
	}

	return changeType, nil
}

// normalizeResponse removes any changes in response due to cache busters added by user for comparison
func normalizeResponse(resp client.Response) client.Response {
	resp.Location = stripInjectedParam(resp.Location, "cache")
	resp.Body = stripParamFromBody(resp.Body, "cache")
	return resp
}

// Removes `cache=xxxxx` from URLs using URL parsing
func stripInjectedParam(rawURL, param string) string {
	if rawURL == "" {
		return ""
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL // fallback
	}

	q := parsed.Query()
	q.Del(param)
	parsed.RawQuery = q.Encode()

	// Remove trailing "?" if empty query
	return strings.TrimSuffix(parsed.String(), "?")
}

// Optional: remove from body if URL is embedded (basic)
func stripParamFromBody(body string, param string) string {
	re := regexp.MustCompile(`[?&]` + param + `=[^&\s">]{5}`)
	return re.ReplaceAllString(body, "")
}

// isEmptyResponse checks if a client.Response is an empty struct with default values.
func isEmptyResponse(resp client.Response) bool {
	return resp.StatusCode == 0 // can never be 0 in a valid response
}
