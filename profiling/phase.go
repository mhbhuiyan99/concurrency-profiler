package profiling

import (
	"runtime"
	"time"
)

// PhaseResult represents the profiling result of one execution phase.
//
// Responsibilities:
//   - Store the phase name.
//   - Store the execution time.
//   - Store memory statistics before and after execution.
//   - Store goroutine counts before and after execution.
type PhaseResult struct {
	Name              string
	ExecutionTime     time.Duration
	MemoryBefore      MemoryStats
	MemoryAfter       MemoryStats
	GoroutinesBefore  int
	GoroutinesAfter   int
}

// ProfilePhase measures the performance of one execution phase.
//
// Responsibilities:
//   - Capture memory statistics before execution.
//   - Capture goroutine count before execution.
//   - Measure execution time.
//   - Execute the provided function.
//   - Capture memory statistics after execution.
//   - Capture goroutine count after execution.
//   - Return the collected profiling metrics.
func ProfilePhase(
	name string,
	fn func(),
) PhaseResult {

	memoryBefore := GetMemoryStats()
	goroutinesBefore := runtime.NumGoroutine()

	start := time.Now()

	fn()

	executionTime := time.Since(start)

	memoryAfter := GetMemoryStats()
	goroutinesAfter := runtime.NumGoroutine()

	return PhaseResult{
		Name:             name,
		ExecutionTime:    executionTime,
		MemoryBefore:     memoryBefore,
		MemoryAfter:      memoryAfter,
		GoroutinesBefore: goroutinesBefore,
		GoroutinesAfter:  goroutinesAfter,
	}
}