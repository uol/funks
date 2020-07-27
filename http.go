package funks

import (
	"crypto/tls"
	"net/http"
	"time"
)

/**
* Common functions used by http protocol.
* @author rnojiri
**/

// CreateHTTPClient - creates a new HTTP client (use zero maxConnsPerHost for unlimited connections)
func CreateHTTPClient(timeout time.Duration, insecureSkipVerify bool, maxConnsPerHost int) *http.Client {

	transportCore := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
		MaxConnsPerHost: maxConnsPerHost,
	}

	httpClient := &http.Client{
		Transport: transportCore,
		Timeout:   timeout,
	}

	return httpClient
}
