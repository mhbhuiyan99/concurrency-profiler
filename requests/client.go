package requests

import (
	"net/http"
	"time"
)

// httpClient is the shared HTTP client used for all external API requests.
//
// Responsibilities:
//   - Reuse TCP connections.
//   - Apply a timeout to every request.
//   - Avoid creating a new client for every API call.
var httpClient = &http.Client {
	Timeout: 10*time.Second,
}
