package routers

import (
	"socket/controllers"
	"github.com/astaxie/beego"
)

func init() {

	beego.Router("/", &controllers.SocketController{}, "*:ToConn")
}
