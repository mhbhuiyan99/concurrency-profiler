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

// TestConcurrency executes the concurrency benchmark and coordinates profiling.
//
// Responsibilities:
//   - Start and stop CPU profiling.
//   - Execute the sequential implementation.
//   - Execute the WaitGroup implementation.
//   - Execute the channel implementation.
//   - Collect controller-level memory statistics.
//   - Return the benchmark results as a JSON response.
func (c *ConcurrencyController) TestConcurrency() {

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