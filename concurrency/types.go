package concurrency

import "time"

// APIResult represents the outcome of a single API request.
//
// Responsibilities:
//   - Store the request URL.
//   - Store the HTTP status code.
//   - Store the response body.
//   - Store any request error.
type APIResult struct {
	URL        string
	StatusCode int
	Response   string
	Err        error
}

// PhaseResult represents the execution result of one concurrency strategy.
//
// Responsibilities:
//   - Store the execution method.
//   - Store the execution time.
//   - Store all API request results.
type PhaseResult struct {
	Name          string
	ExecutionTime time.Duration
	Results       []APIResult
}
