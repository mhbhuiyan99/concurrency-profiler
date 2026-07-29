package controllers

import (
	"concurrency-profiler/concurrency"
	"concurrency-profiler/profiling"
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
	fmt.Printf("After: %+v\n", after)

	c.Data["json"] = map[string]any{
		"sequential": seq,
		"waitgroup":  wg,
		"channel":    ch,
	}

	c.ServeJSON()
}