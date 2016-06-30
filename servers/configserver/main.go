package main

import (
	_ "configserver/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
