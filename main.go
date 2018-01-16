package main

import (
	_ "socket/routers"
	"github.com/astaxie/beego"
	"flag"
)

func main() {

	port := flag.Int("port", 9105, "http port")
	flag.Parse()

	beego.BConfig.Listen.HTTPPort = *port
	beego.Run()
}
