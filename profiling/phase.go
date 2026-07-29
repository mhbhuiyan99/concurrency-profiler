package profiling

import (
	"concurrency-profiler/concurrency"
	"concurrency-profiler/utils"
)

// PhaseResult represents profiling information for one execution phase.
//
// Responsibilities:
//   - Store the execution result.
//   - Store memory statistics before and after execution.
type PhaseResult struct {
	Execution   concurrency.PhaseResult
	MemoryBefore MemoryStats
	MemoryAfter  MemoryStats
}

// ProfilePhase profiles one execution phase.
//
// Responsibilities:
//   - Capture memory statistics before execution.
//   - Execute the provided function using the common execution timer.
//   - Capture memory statistics after execution.
//   - Return the execution and profiling results.
func ProfilePhase(
	name string,
	fn func() []concurrency.APIResult,
) PhaseResult {

	memoryBefore := GetMemoryStats()
	
	execution := utils.MeasureExecution(
		name,
		fn,
	)

	memoryAfter := GetMemoryStats()

	return PhaseResult{
		Execution:        execution,
		MemoryBefore:     memoryBefore,
		MemoryAfter:      memoryAfter,
	}
}