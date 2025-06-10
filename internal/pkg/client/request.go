// Description: This file contains functions for making HTTP requests
// including fetching responses and sending requests without reading responses.

package client

import (
	"fmt"
	"io"
	"net/http"
)

// Response holds the response from the URL
type Response struct {
	StatusCode int                 `json:"StatusCode"`
	Headers    map[string][]string `json:"Headers"`
	Body       string              `json:"Body"`
	Location   string              `json:"Location"`
}

// FetchResponse sends a GET request and returns the response.
func FetchResponse(url string, requestHeaders map[string]string, httpClient *http.Client) (*Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &Response{}, fmt.Errorf("error creating HTTP request: %v", err)
	}

	// Add headers to request
	for key, value := range requestHeaders {
		req.Header.Add(key, value)
	}

	// Send request
	resp, err := httpClient.Do(req)
	if err != nil {
		return &Response{}, fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Response{}, fmt.Errorf("error reading response body: %v", err)
	}

	// Return response
	return &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       string(respBody),
		Location:   resp.Header.Get("Location"),
	}, nil
}

// SendRequest sends a request but does not read or return the response.
func SendRequest(url string, requestHeaders map[string]string, httpClient *http.Client) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	// Add headers to request
	for key, value := range requestHeaders {
		req.Header.Add(key, value)
	}

	// Send request (ignoring the response)
	_, err = httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %v", err)
	}

	return nil
}
