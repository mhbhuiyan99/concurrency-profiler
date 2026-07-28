package concurrency

import (
	"concurrency-profiler/requests"
	"time"
)

// RunSequential executes all API requests one after another.
//
// Responsibilities:
//   - Execute each API request sequentially.
//   - Collect the result of every request.
//   - Measure the total execution time.
//   - Return the execution summary.
func RunSequential(urls []string) PhaseResult {
	start := time.Now()

	results := make([]APIResult, 0, len(urls))

	for _, url := range urls {
		req, err := requests.NewGetRequest(url)
		if err != nil {
			results = append(results, APIResult{
				URL: url,
				Err: err,
			})
			continue
		}

		body, statusCode, err := requests.DoRequest(req)

		results = append(results, APIResult{
			URL:        url,
			StatusCode: statusCode,
			Response:   body,
			Err:        err,
		})
	}

	return PhaseResult{
		Name:          "Sequential",
		ExecutionTime: time.Since(start),
		Results:       results,
	}
}
