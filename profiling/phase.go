package profiling

import (
	"concurrency-profiler/concurrency"
	"concurrency-profiler/utils"
	"runtime"
)

// PhaseResult represents profiling information for one execution phase.
//
// Responsibilities:
//   - Store the execution result.
//   - Store memory statistics before and after execution.
//   - Store goroutine counts before and after execution.
type PhaseResult struct {
	Execution         concurrency.PhaseResult
	MemoryBefore      MemoryStats
	MemoryAfter       MemoryStats
	GoroutinesBefore  int
	GoroutinesAfter   int
}

// ProfilePhase profiles one execution phase.
//
// Responsibilities:
//   - Capture memory statistics before execution.
//   - Capture goroutine count before execution.
//   - Execute the provided function using the common execution timer.
//   - Capture memory statistics after execution.
//   - Capture goroutine count after execution.
//   - Return the execution and profiling results.
func ProfilePhase(
	name string,
	fn func() []concurrency.APIResult,
) PhaseResult {

	memoryBefore := GetMemoryStats()
	goroutinesBefore := runtime.NumGoroutine()

	execution := utils.MeasureExecution(name, fn)

	memoryAfter := GetMemoryStats()
	goroutinesAfter := runtime.NumGoroutine()

	return PhaseResult{
		Execution:        execution,
		MemoryBefore:     memoryBefore,
		MemoryAfter:      memoryAfter,
		GoroutinesBefore: goroutinesBefore,
		GoroutinesAfter:  goroutinesAfter,
	}
}