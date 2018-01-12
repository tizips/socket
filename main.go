package main

import (
	_ "socket/routers"

	"github.com/astaxie/beego"
)

func main() {

	beego.Run()
}
