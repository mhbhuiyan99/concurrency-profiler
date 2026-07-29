package concurrency

import "concurrency-profiler/requests"

// FetchAPI executes a single API request.
//
// Responsibilities:
//   - Create the HTTP request.
//   - Execute the request.
//   - Collect the response.
//   - Return the request result.
func FetchAPI(url string) APIResult {

	req, err := requests.NewGETRequest(url)
	if err != nil {
		return APIResult{
			URL: url,
			Err: err,
		}
	}

	body, statusCode, err := requests.DoRequest(req)

	return APIResult{
		URL:        url,
		StatusCode: statusCode,
		Response:   string(body),
		Err:        err,
	}
}