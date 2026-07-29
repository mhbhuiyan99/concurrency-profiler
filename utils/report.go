package utils

import (
	"fmt"
	"time"
)

// MemoryStatsView contains memory statistics required by the terminal report.
//
// Responsibilities:
//   - Provide memory statistics to the reporting layer.
//   - Keep the reporting layer independent from the profiling package.
type MemoryStatsView struct {
	Alloc      uint64
	TotalAlloc uint64
	Sys        uint64
	NumGC      uint32
}

// PrintPerformanceHeader prints the performance test header.
//
// Responsibilities:
//   - Display the report title.
//   - Display the total number of APIs being tested.
func PrintPerformanceHeader(totalAPIs int) {
	fmt.Println()
	fmt.Println("================ API PERFORMANCE TEST ================")
	fmt.Println()
	fmt.Printf("Total APIs Called: %d\n", totalAPIs)
	fmt.Println()
}

// PrintExecutionTime prints the execution time of one execution method.
//
// Responsibilities:
//   - Display the execution method.
//   - Display the time taken by the method.
func PrintExecutionTime(
	name string,
	executionTime time.Duration,
) {
	fmt.Printf("[ %s Execution ]\n", name)
	fmt.Println("-------------------------------------")
	fmt.Printf("Time Taken : %.2f ms\n", executionTime.Seconds()*1000)
	fmt.Println()
}

// PrintMemoryStats prints controller-level memory statistics.
//
// Responsibilities:
//   - Display memory statistics before execution.
//   - Display memory statistics after execution.
//   - Display Alloc, TotalAlloc, Sys, and NumGC values.
func PrintMemoryStats(
	before MemoryStatsView,
	after MemoryStatsView,
) {
	fmt.Println("--------------- MEMORY STATS ----------------")
	fmt.Println()

	fmt.Println("Before Execution")
	fmt.Printf("Alloc      : %.2f MB\n", BytesToMB(before.Alloc))
	fmt.Printf("TotalAlloc : %.2f MB\n", BytesToMB(before.TotalAlloc))
	fmt.Printf("Sys        : %.2f MB\n", BytesToMB(before.Sys))
	fmt.Printf("NumGC      : %d\n", before.NumGC)
	fmt.Println()

	fmt.Println("After Execution")
	fmt.Printf("Alloc      : %.2f MB\n", BytesToMB(after.Alloc))
	fmt.Printf("TotalAlloc : %.2f MB\n", BytesToMB(after.TotalAlloc))
	fmt.Printf("Sys        : %.2f MB\n", BytesToMB(after.Sys))
	fmt.Printf("NumGC      : %d\n", after.NumGC)
	fmt.Println()
}

// PrintPhaseProfiling prints profiling information for each execution phase.
//
// Responsibilities:
//   - Display execution time for each phase.
//   - Display allocation volume for each phase.
//
// Note:
//   - Bytes Allocated represents allocation volume during the phase.
//   - Goroutine profiling is intentionally omitted based on instructor guidance.
func PrintPhaseProfiling(
	name string,
	executionTime time.Duration,
	bytesAllocated uint64,
) {
	fmt.Println(name)
	fmt.Printf("Time Taken      : %.2f ms\n", executionTime.Seconds()*1000)
	fmt.Printf("Bytes Allocated : +%.2f MB\n", BytesToMB(bytesAllocated))
	fmt.Println()
}

// PrintPerformanceComparison prints the performance comparison.
//
// Responsibilities:
//   - Display the comparison section.
//   - Display visual execution-time bars.
//   - Display performance gains between execution methods.
//   - Display allocation volume comparison.
func PrintPerformanceComparison(
	sequentialTime time.Duration,
	waitGroupTime time.Duration,
	channelTime time.Duration,
	waitGroupGain float64,
	channelGain float64,
	channelVsWaitGroup float64,
	highestAllocationMethod string,
	lowestAllocationMethod string,
) {
	fmt.Println("================= COMPARISON =================")
	fmt.Println()

	PrintPerformanceBars(
		sequentialTime,
		waitGroupTime,
		channelTime,
	)

	fmt.Println("Performance Gain")
	fmt.Println("-------------------------------------")

	fmt.Printf(
		"WaitGroup vs Sequential : %s\n",
		FormatPerformanceGain(waitGroupGain),
	)

	fmt.Printf(
		"Channel vs Sequential   : %s\n",
		FormatPerformanceGain(channelGain),
	)

	fmt.Printf(
		"Channel vs WaitGroup    : %s\n",
		FormatPerformanceGain(channelVsWaitGroup),
	)

	fmt.Println()

	fmt.Println("Allocation Volume")
	fmt.Println("-------------------------------------")
	fmt.Printf("Highest Allocation : %s\n", highestAllocationMethod)
	fmt.Printf("Lowest Allocation  : %s\n", lowestAllocationMethod)
	fmt.Println()
}

// PrintProfilingReport prints the controller-level profiling report.
//
// Responsibilities:
//   - Display total controller execution time.
//   - Display controller-level memory statistics.
//   - Display phase-level profiling results.
//   - Display CPU profiling status.
func PrintProfilingReport(
	totalExecutionTime time.Duration,
	before MemoryStatsView,
	after MemoryStatsView,
	sequentialTime time.Duration,
	sequentialBytes uint64,
	waitGroupTime time.Duration,
	waitGroupBytes uint64,
	channelTime time.Duration,
	channelBytes uint64,
) {
	fmt.Println("================ PROFILING REPORT ================")
	fmt.Println()

	fmt.Println("Controller : /test-concurrency")
	fmt.Println()
	fmt.Printf(
		"Total Execution Time : %.2f ms\n",
		totalExecutionTime.Seconds()*1000,
	)
	fmt.Println()

	PrintMemoryStats(before, after)

	fmt.Println("--------------- PHASE PROFILING ----------------")
	fmt.Println()

	PrintPhaseProfiling(
		"Sequential",
		sequentialTime,
		sequentialBytes,
	)

	PrintPhaseProfiling(
		"WaitGroup",
		waitGroupTime,
		waitGroupBytes,
	)

	PrintPhaseProfiling(
		"Channel",
		channelTime,
		channelBytes,
	)

	fmt.Println("--------------- CPU PROFILE ----------------")
	fmt.Println()
	fmt.Println("CPU Profiling Started...")
	fmt.Println("CPU Profiling Stopped.")
	fmt.Println()
}

// PrintSummary prints the final performance analysis.
//
// Responsibilities:
//   - Display the fastest execution method.
//   - Display the highest allocation method.
//   - Display the most efficient execution method.
func PrintSummary(
	fastestMethod string,
	highestAllocationMethod string,
	mostEfficientMethod string,
) {
	fmt.Println("--------------- SUMMARY ----------------")
	fmt.Println()

	fmt.Printf("Fastest Method           : %s\n", fastestMethod)
	fmt.Printf("Highest Allocation       : %s\n", highestAllocationMethod)
	fmt.Printf("Most Efficient Method    : %s\n", mostEfficientMethod)

	fmt.Println()
	fmt.Println("====================================================")
}

// FormatPerformanceGain formats a performance comparison result.
//
// Responsibilities:
//   - Describe positive values as faster.
//   - Describe negative values as slower.
//   - Return a readable performance comparison.
func FormatPerformanceGain(value float64) string {
	if value >= 0 {
		return fmt.Sprintf("%.2f%% faster", value)
	}

	return fmt.Sprintf("%.2f%% slower", -value)
}

// BytesToMB converts bytes to megabytes.
//
// Responsibilities:
//   - Convert byte values into megabytes.
//   - Keep unit conversion reusable across reports.
func BytesToMB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024)
}