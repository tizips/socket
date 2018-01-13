package routers

import (
	"socket/controllers"
	"github.com/astaxie/beego"
)

func init() {

	beego.Router("/", &controllers.SocketController{}, "*:ToConn")

	beego.Router("/send", &controllers.SocketController{}, "*:ToSend")
}
