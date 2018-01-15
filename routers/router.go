package routers

import (
	"socket/controllers"
	"github.com/astaxie/beego"
)

func init() {

	beego.Router("/ws", &controllers.SocketController{}, "*:ToConn")

	beego.Router("/send", &controllers.SocketController{}, "*:ToSend")
}
