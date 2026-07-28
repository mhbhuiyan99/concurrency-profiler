package utils

import (
	"time"

	"concurrency-profiler/concurrency"
)

// MeasureExecution measures how long an execution strategy takes.
//
// Responsibilities:
//   - Record the start time.
//   - Execute the provided function.
//   - Calculate the elapsed time.
//   - Return a completed PhaseResult.
func MeasureExecution(
	name string,
	fn func() []concurrency.APIResult,
) concurrency.PhaseResult {

	start := time.Now()

	results := fn()

	return concurrency.PhaseResult{
		Name:          name,
		ExecutionTime: time.Since(start),
		Results:       results,
	}
}