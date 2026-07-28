package controllers

import (
	"concurrency-profiler/concurrency"
	"concurrency-profiler/utils"

	"github.com/beego/beego/v2/server/web"
)

type ConcurrencyController struct {
	web.Controller
}

// TestConcurrency executes and compares all concurrency implementations.
//
// Responsibilities:
//   - Execute the sequential implementation.
//   - Execute the WaitGroup implementation.
//   - Execute the channel implementation.
//   - Return the execution results as a JSON response.
func (c *ConcurrencyController) TestConcurrency() {

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

	c.Data["json"] = map[string]any{
		"sequential": seq,
		"waitgroup":  wg,
		"channel":    ch,
	}
	c.ServeJSON()
}