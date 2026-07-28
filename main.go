package main

import (
	_ "concurrency-profiler/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}
