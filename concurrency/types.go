package concurrency

import "time"

// APIResult represents the outcome of a single API request.
//
// Responsibilities:
//   - Store the request URL.
//   - Store the response body.
//   - Store any request error.
type APIResult struct {
	URL      string
	Response []byte
	Err      error
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