package routers

import (
	"concurrency-profiler/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/test-concurrency", &controllers.ConcurrencyController{}, "get:TestConcurrency")
}