package controllers

import "github.com/beego/beego/v2/server/web"

type ConcurrencyController struct {
	web.Controller
}

func (c *ConcurrencyController) TestConcurrency() {
	c.Data["json"] = map[string]string {
		"message": "Concurrency Profile API is running",
	}
	c.ServeJSON()
}