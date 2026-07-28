package controllers

import (
	"concurrency-profiler/concurrency"
	"concurrency-profiler/utils"

	"github.com/beego/beego/v2/server/web"
)

type ConcurrencyController struct {
	web.Controller
}

// TestConcurrency executes the sequential concurrency benchmark.
//
// Responsibilities:
//   - Execute the sequential implementation.
//   - Return the benchmark result as JSON.
func (c *ConcurrencyController) TestConcurrency() {

	result := utils.MeasureExecution(
		"Sequential",
		func() []concurrency.APIResult {
			return concurrency.RunSequential(concurrency.APIURLs)
		},
	)

	c.Data["json"] = result
	c.ServeJSON()
}