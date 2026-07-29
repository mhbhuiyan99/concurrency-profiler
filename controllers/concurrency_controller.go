package controllers

import (
	"concurrency-profiler/concurrency"
	"concurrency-profiler/profiling"
	"concurrency-profiler/utils"
	"fmt"
	"time"

	"github.com/beego/beego/v2/server/web"
)

type ConcurrencyController struct {
	web.Controller
}

// TestConcurrency executes the concurrency benchmark and coordinates profiling.
//
// Responsibilities:
//   - Start and stop CPU profiling.
//   - Execute the sequential implementation.
//   - Execute the WaitGroup implementation.
//   - Execute the channel implementation.
//   - Collect controller-level memory statistics.
//   - Generate the terminal performance report.
//   - Return the benchmark results as a JSON response.
func (c *ConcurrencyController) TestConcurrency() {

	start := time.Now()

	cpuProfile, err := profiling.StartCPUProfile("profiles/cpu.prof")
	if err != nil {
		c.Data["json"] = map[string]any{
			"error": err.Error(),
		}
		c.ServeJSON()
		return
	}

	defer func() {
		if err := profiling.StopCPUProfile(cpuProfile); err != nil {
			fmt.Printf("CPU profiling stop failed: %v\n", err)
		}
	}()

	before := profiling.GetMemoryStats()

	seq := profiling.ProfilePhase(
		"Sequential",
		func() []concurrency.APIResult {
			return concurrency.RunSequential(concurrency.APIURLs)
		},
	)

	wg := profiling.ProfilePhase(
		"WaitGroup",
		func() []concurrency.APIResult {
			return concurrency.RunWaitGroup(concurrency.APIURLs)
		},
	)

	ch := profiling.ProfilePhase(
		"Channel",
		func() []concurrency.APIResult {
			return concurrency.RunChannel(concurrency.APIURLs)
		},
	)

	after := profiling.GetMemoryStats()

	comparison := profiling.ComparePhases(
		seq,
		wg,
		ch,
	)

	totalExecutionTime := time.Since(start)

	// Print the terminal performance report.
	utils.PrintPerformanceHeader(len(concurrency.APIURLs))

	utils.PrintExecutionTime(
		"Sequential",
		seq.Execution.ExecutionTime,
	)

	utils.PrintExecutionTime(
		"WaitGroup",
		wg.Execution.ExecutionTime,
	)

	utils.PrintExecutionTime(
		"Channel",
		ch.Execution.ExecutionTime,
	)

	utils.PrintPerformanceComparison(
		seq.Execution.ExecutionTime,
		wg.Execution.ExecutionTime,
		ch.Execution.ExecutionTime,
		comparison.WaitGroupGain,
		comparison.ChannelGain,
		comparison.WaitGroupVsChannel,
		comparison.HighestMemoryMethod,
		comparison.LowestMemoryMethod,
	)

	utils.PrintProfilingReport(
		totalExecutionTime,

		utils.MemoryStatsView{
			Alloc:      before.Alloc,
			TotalAlloc: before.TotalAlloc,
			Sys:        before.Sys,
			NumGC:      before.NumGC,
		},

		utils.MemoryStatsView{
			Alloc:      after.Alloc,
			TotalAlloc: after.TotalAlloc,
			Sys:        after.Sys,
			NumGC:      after.NumGC,
		},

		seq.Execution.ExecutionTime,
		profiling.CalculateMemoryUsed(seq),

		wg.Execution.ExecutionTime,
		profiling.CalculateMemoryUsed(wg),

		ch.Execution.ExecutionTime,
		profiling.CalculateMemoryUsed(ch),
	)

	utils.PrintSummary(
		comparison.FastestMethod,
		comparison.HighestMemoryMethod,
		comparison.MostEfficientMethod,
	)

	c.Data["json"] = map[string]any{
		"sequential":         seq,
		"waitgroup":          wg,
		"channel":            ch,
		"comparison":         comparison,
		"totalExecutionTime": totalExecutionTime,
	}

	c.ServeJSON()
}
