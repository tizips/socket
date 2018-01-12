package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"fmt"
)

type SocketController struct {
	beego.Controller
}

type Response struct {
	Msg    string            `json:"msg"`
	Code   int               `json:"code"`
	Result map[string]string `json:"result"`
}

type Client struct {
	socket *websocket.Conn
	typing string
	hotel  string
}

var (
	response Response
	upgrader = websocket.Upgrader{}
)

func (this *SocketController) Prepare() {
	response.Code = 0
	response.Msg = "success"
	response.Result = nil
}

func (this *SocketController) ToConn() {

	ws, err := upgrader.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil)

	if err != nil {
		fmt.Println(err)
	}
	//clients[ws] = true

	//不断的广播发送到页面上

	if err := ws.WriteJSON(response); err != nil {
		fmt.Println(err)
	}
}
