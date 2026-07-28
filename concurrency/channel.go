package concurrency

import "sync"

// RunChannel executes all API requests concurrently using
// goroutines and channels.
//
// Responsibilities:
//   - Launch one goroutine per API request.
//   - Send each API result through a channel.
//   - Collect all results from the channel.
//   - Return all request results.
func RunChannel(urls []string) []APIResult {
	results := make([]APIResult, 0, len(urls))

	// Using Buffered channel
	// Because we know exactly how many API calls we'll make
	resultsChan := make(chan APIResult, len(urls))

	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)

		go func(endpoint string) {
			defer wg.Done()
			resultsChan <- FetchAPI(endpoint)
		}(url)
	}

	// Close the channel after all goroutines finish
	go func() {
		wg.Wait()
		close(resultsChan)
	}()


	for result := range resultsChan {
		results = append(results, result)
	}

	return results
}