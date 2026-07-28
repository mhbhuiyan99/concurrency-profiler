package concurrency

import (
	"sync"
)

// RunWaitGroup executes all API requests concurrently using goroutines
// synchronized with a sync.WaitGroup.
//
// Responsibilities:
//   - Launch one goroutine per API request.
//   - Wait for all goroutines to finish.
//   - Collect all request results.
//   - Return all request results.
func RunWaitGroup(urls []string) []APIResult {
	results := make([]APIResult, len(urls))

	var wg sync.WaitGroup

	for i, url := range urls {
		wg.Add(1)

		go func(index int, endpoint string) {
			defer wg.Done()

			results[index] = FetchAPI(endpoint)
		} (i, url)
	}

	wg.Wait()

	return results
}