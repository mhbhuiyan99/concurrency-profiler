// Package requests contains reusable HTTP request helpers.
//
// Responsibilities:
//   - Create configured HTTP requests.
//   - Apply authentication and default headers.
//   - Execute external API requests.
//   - Return response bodies or descriptive errors.
package requests

import (
	"fmt"
	"io"
	"net/http"

	"github.com/beego/beego/v2/server/web"
)

// NewGETRequest creates a configured HTTP GET request.
//
// Responsibilities:
//   - Create a new HTTP GET request.
//   - Apply the application's default authentication.
//   - Apply the application's default request headers.
//   - Return the configured request.
func NewGETRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	if err := setDefaultHeaders(req); err != nil {
		return nil, fmt.Errorf("set default headers failed: %w", err)
	}

	return req, nil
}

// DoRequest executes an HTTP request and returns the response body.
//
// Responsibilities:
//   - Execute the HTTP request.
//   - Validate the HTTP status code.
//   - Read the response body.
//   - Return the response body, status code, and descriptive errors.
func DoRequest(req *http.Request) ([]byte, int, error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return nil, statusCode, fmt.Errorf("unexpected HTTP status: %d", statusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, statusCode, fmt.Errorf("read response body failed: %w", err)
	}

	return body, statusCode, nil
}

// setDefaultHeaders applies the application's default
// authentication and HTTP headers to a request.
//
// Responsibilities:
//   - Apply Basic Authentication.
//   - Set common request headers.
//   - Set the API key.
//   - Prepare the request for external API communication.
func setDefaultHeaders(req *http.Request) error {
	username, err := web.AppConfig.String("username")
	if err != nil {
		return err
	}

	password, err := web.AppConfig.String("password")
	if err != nil {
		return err
	}

	apiKey, err := web.AppConfig.String("api_key")
	if err != nil {
		return err
	}

	req.SetBasicAuth(username, password)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", "desktop")
	req.Header.Set("Origin", "123presto-MS-ROW.com")
	req.Header.Set("x-api-key", apiKey)

	return nil
}