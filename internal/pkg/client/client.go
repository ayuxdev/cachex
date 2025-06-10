// Description: This file contains the client package which is responsible for
// creating a new HTTP client with custom transport settings.
package client

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
)

// Config holds configuration settings for the HTTP client.
type Config struct {
	DialTimeout           time.Duration // Timeout for establishing the connection
	HandshakeTimeout      time.Duration // Timeout for TLS handshake
	ResponseHeaderTimeout time.Duration // Timeout for server response headers
	ProxyURL              string        // Proxy URL for the HTTP client (optional)
}

// New creates a new HTTP client with custom transport settings.
func (c Config) CreateNewClient() *http.Client {
	// Custom transport settings for the HTTP client.
	CustomTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Ignore invalid SSL certificates
		DialContext: (&net.Dialer{
			Timeout: c.DialTimeout, // Timeout for establishing the connection
		}).DialContext,
		TLSHandshakeTimeout:   c.HandshakeTimeout,      // Timeout for TLS handshake
		ResponseHeaderTimeout: c.ResponseHeaderTimeout, // Timeout for server response headers
	}

	if c.ProxyURL != "" {
		proxyURL, err := url.Parse(c.ProxyURL)
		if err == nil {
			CustomTransport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	return &http.Client{Transport: CustomTransport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Prevents following redirects
		}}
}
