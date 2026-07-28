package concurrency

import (
	"concurrency-profiler/requests"
)

// RunSequential executes all API requests one after another.
//
// Responsibilities:
//   - Execute each API request sequentially.
//   - Collect the result of every request.
//   - Return all request results.
func RunSequential(urls []string) []APIResult {

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
			Response:   string(body),
			Err:        err,
		})
	}

	return results
}
