package controllers

import (
	"concurrency-profiler/concurrency"
	"concurrency-profiler/profiling"
	"concurrency-profiler/utils"
	"fmt"

	"github.com/beego/beego/v2/server/web"
)

type ConcurrencyController struct {
	web.Controller
}

// TestConcurrency executes the concurrency benchmark.
//
// Responsibilities:
//   - Execute the sequential implementation.
//   - Execute the WaitGroup implementation.
//   - Execute the channel implementation.
//   - Collect controller-level memory statistics.
//   - Return the benchmark results as a JSON response.
func (c *ConcurrencyController) TestConcurrency() {

	before := profiling.GetMemoryStats()
	fmt.Printf("Before: %+v\n", before)

	seq := utils.MeasureExecution(
		"Sequential",
		func() []concurrency.APIResult {
			return concurrency.RunSequential(concurrency.APIURLs)
		},
	)

	wg := utils.MeasureExecution(
		"WaitGroup",
		func() []concurrency.APIResult {
			return concurrency.RunWaitGroup(concurrency.APIURLs)
		},
	)

	ch := utils.MeasureExecution(
		"Channel",
		func() []concurrency.APIResult {
			return concurrency.RunChannel(concurrency.APIURLs)
		},
	)

	after := profiling.GetMemoryStats()
	fmt.Printf("After: %+v\n", after)

	c.Data["json"] = map[string]any{
		"sequential": seq,
		"waitgroup":  wg,
		"channel":    ch,
	}
	c.ServeJSON()
}