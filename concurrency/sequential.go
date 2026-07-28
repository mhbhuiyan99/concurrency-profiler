package concurrency

// RunSequential executes all API requests one after another.
//
// Responsibilities:
//   - Execute each API request sequentially.
//   - Collect the result of every request.
//   - Return all request results.
func RunSequential(urls []string) []APIResult {

	results := make([]APIResult, 0, len(urls))

	for _, url := range urls {
		results = append(results, FetchAPI(url))
	}

	return results
}