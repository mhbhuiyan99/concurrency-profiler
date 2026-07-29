package utils

import (
	"fmt"
	"time"
)

const (
	visualBarWidth = 30
	visualBarChar  = "█"
)

// PrintPerformanceBars prints a visual comparison of execution times.
//
// Responsibilities:
//   - Compare execution times of all execution methods.
//   - Scale each method relative to the slowest method.
//   - Display the comparison using terminal bars.
func PrintPerformanceBars(
	sequential time.Duration,
	waitGroup time.Duration,
	channel time.Duration,
) {
	fmt.Println("Execution Time")
	fmt.Println("-------------------------------------")

	maxTime := maxDuration(
		sequential,
		waitGroup,
		channel,
	)

	printPerformanceBar("Sequential", sequential, maxTime)
	printPerformanceBar("WaitGroup", waitGroup, maxTime)
	printPerformanceBar("Channel", channel, maxTime)

	fmt.Println()
}

// printPerformanceBar prints one scaled performance bar.
//
// Responsibilities:
//   - Calculate the relative bar length.
//   - Display the execution method and its bar.
func printPerformanceBar(
	name string,
	executionTime time.Duration,
	maxTime time.Duration,
) {
	barLength := calculateBarLength(executionTime, maxTime)

	fmt.Printf(
		"%-10s | %-30s %.2f s\n",
		name,
		repeatBar(barLength),
		executionTime.Seconds(),
	)
}

// calculateBarLength calculates the visual bar length relative to
// the slowest execution method.
//
// Responsibilities:
//   - Scale the execution time against the maximum execution time.
//   - Return a bar length between 1 and visualBarWidth.
func calculateBarLength(
	executionTime time.Duration,
	maxTime time.Duration,
) int {
	if executionTime <= 0 || maxTime <= 0 {
		return 0
	}

	length := int(
		float64(executionTime) / float64(maxTime) *
			float64(visualBarWidth),
	)

	if length < 1 {
		length = 1
	}

	return length
}

// repeatBar creates a terminal bar of the requested length.
//
// Responsibilities:
//   - Generate a repeated visual character.
//   - Return the terminal representation of the performance bar.
func repeatBar(length int) string {
	result := ""

	for i := 0; i < length; i++ {
		result += visualBarChar
	}

	return result
}

// maxDuration returns the largest execution duration.
//
// Responsibilities:
//   - Compare all execution durations.
//   - Return the largest duration.
func maxDuration(
	sequential time.Duration,
	waitGroup time.Duration,
	channel time.Duration,
) time.Duration {
	max := sequential

	if waitGroup > max {
		max = waitGroup
	}

	if channel > max {
		max = channel
	}

	return max
}
