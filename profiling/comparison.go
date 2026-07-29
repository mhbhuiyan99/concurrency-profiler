package profiling

import (
	"fmt"
	"time"
)

// ComparisonResult contains the performance comparison of all execution phases.
//
// Responsibilities:
//   - Identify the fastest execution method.
//   - Calculate performance gains between execution methods.
//   - Compare allocation volume between execution methods.
//   - Identify the overall most efficient method.
type ComparisonResult struct {
	FastestMethod       string
	SequentialGain      float64
	WaitGroupGain       float64
	ChannelGain         float64
	WaitGroupVsChannel  float64
	HighestMemoryMethod string
	LowestMemoryMethod  string
	SequentialMemory    uint64
	WaitGroupMemory     uint64
	ChannelMemory       uint64
	MostEfficientMethod string
}

// CalculateMemoryUsed returns the amount of memory allocated during a phase.
//
// Responsibilities:
//   - Compare TotalAlloc before and after execution.
//   - Return the difference as the phase memory usage.
func CalculateMemoryUsed(result PhaseResult) uint64 {
	if result.MemoryAfter.TotalAlloc < result.MemoryBefore.TotalAlloc {
		return 0
	}

	return result.MemoryAfter.TotalAlloc - result.MemoryBefore.TotalAlloc
}

// CalculatePerformanceGain calculates how much faster a method is
// compared with the baseline execution time.
//
// Responsibilities:
//   - Compare the baseline execution time with the current execution time.
//   - Return the performance gain as a percentage.
func CalculatePerformanceGain(
	baseline time.Duration,
	current time.Duration,
) float64 {

	if baseline <= 0 {
		return 0
	}

	return float64(baseline-current) / float64(baseline) * 100
}

// ComparePhases compares the performance of all execution methods.
//
// Responsibilities:
//   - Identify the fastest execution method.
//   - Calculate performance gains.
//   - Compare allocation volume between execution methods.
//   - Identify the overall most efficient method.
func ComparePhases(
	sequential PhaseResult,
	waitGroup PhaseResult,
	channel PhaseResult,
) ComparisonResult {

	sequentialTime := sequential.Execution.ExecutionTime
	waitGroupTime := waitGroup.Execution.ExecutionTime
	channelTime := channel.Execution.ExecutionTime

	sequentialMemory := CalculateMemoryUsed(sequential)
	waitGroupMemory := CalculateMemoryUsed(waitGroup)
	channelMemory := CalculateMemoryUsed(channel)

	result := ComparisonResult{
		SequentialGain:      0,
		WaitGroupGain:       CalculatePerformanceGain(sequentialTime, waitGroupTime),
		ChannelGain:         CalculatePerformanceGain(sequentialTime, channelTime),
		WaitGroupVsChannel:  CalculatePerformanceGain(waitGroupTime, channelTime),
		SequentialMemory:    sequentialMemory,
		WaitGroupMemory:     waitGroupMemory,
		ChannelMemory:       channelMemory,
	}

	result.FastestMethod = fastestMethod(
		sequentialTime,
		waitGroupTime,
		channelTime,
	)

	result.HighestMemoryMethod = highestMemoryMethod(
		sequentialMemory,
		waitGroupMemory,
		channelMemory,
	)

	result.LowestMemoryMethod = lowestMemoryMethod(
		sequentialMemory,
		waitGroupMemory,
		channelMemory,
	)

	result.MostEfficientMethod = mostEfficientMethod(
		sequentialTime,
		waitGroupTime,
		channelTime,
		sequentialMemory,
		waitGroupMemory,
		channelMemory,
		result.FastestMethod,
	)

	return result
}

// fastestMethod identifies the execution method with the shortest time.
//
// Responsibilities:
//   - Compare the execution times of all phases.
//   - Return the name of the fastest method.
func fastestMethod(
	sequential time.Duration,
	waitGroup time.Duration,
	channel time.Duration,
) string {

	fastest := sequential
	method := "Sequential"

	if waitGroup < fastest {
		fastest = waitGroup
		method = "WaitGroup"
	}

	if channel < fastest {
		method = "Channel"
	}

	return method
}

// highestMemoryMethod identifies the execution method with the highest
// TotalAlloc usage.
//
// Responsibilities:
//   - Compare memory usage across all phases.
//   - Return the name of the method with the highest usage.
func highestMemoryMethod(
	sequential uint64,
	waitGroup uint64,
	channel uint64,
) string {

	highest := sequential
	method := "Sequential"

	if waitGroup > highest {
		highest = waitGroup
		method = "WaitGroup"
	}

	if channel > highest {
		method = "Channel"
	}

	return method
}

// lowestMemoryMethod identifies the execution method with the lowest
// TotalAlloc usage.
//
// Responsibilities:
//   - Compare memory usage across all phases.
//   - Return the name of the method with the lowest usage.
func lowestMemoryMethod(
	sequential uint64,
	waitGroup uint64,
	channel uint64,
) string {

	lowest := sequential
	method := "Sequential"

	if waitGroup < lowest {
		lowest = waitGroup
		method = "WaitGroup"
	}

	if channel < lowest {
		method = "Channel"
	}

	return method
}

// FormatPercentage formats a performance percentage for terminal output.
//
// Responsibilities:
//   - Convert a performance percentage into readable text.
func FormatPercentage(value float64) string {
	return fmt.Sprintf("%.2f%%", value)
}

// mostEfficientMethod identifies the execution method with the best
// overall performance based on execution time and allocation volume.
//
// Responsibilities:
//   - Prefer a method that is both faster and has lower allocation volume.
//   - Return the fastest method when no method dominates both metrics.
func mostEfficientMethod(
	sequentialTime time.Duration,
	waitGroupTime time.Duration,
	channelTime time.Duration,
	sequentialMemory uint64,
	waitGroupMemory uint64,
	channelMemory uint64,
	fastest string,
) string {
	// Channel dominates both metrics.
	if channelTime <= sequentialTime &&
		channelTime <= waitGroupTime &&
		channelMemory <= sequentialMemory &&
		channelMemory <= waitGroupMemory {
		return "Channel"
	}

	// WaitGroup dominates both metrics.
	if waitGroupTime <= sequentialTime &&
		waitGroupTime <= channelTime &&
		waitGroupMemory <= sequentialMemory &&
		waitGroupMemory <= channelMemory {
		return "WaitGroup"
	}

	// Sequential dominates both metrics.
	if sequentialTime <= waitGroupTime &&
		sequentialTime <= channelTime &&
		sequentialMemory <= waitGroupMemory &&
		sequentialMemory <= channelMemory {
		return "Sequential"
	}

	// No method dominates both metrics.
	return fastest
}