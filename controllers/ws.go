package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"fmt"
	"net/http"
	"encoding/json"
)

type SocketController struct {
	beego.Controller
}

type Response struct {
	Msg    string            `json:"msg"`
	Code   int               `json:"code"`
	Result map[string]string `json:"result"`
}

type Info struct {
	typing  string
	hotel   int
	results map[string]string
}

var (
	response Response
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	info    Info
	clients = make(map[*websocket.Conn]Info)
)

func (this *SocketController) Prepare() {
	response.Code = 0
	response.Msg = "success"
	response.Result = nil
	info.results = nil
	info.results = make(map[string]string)
}

func (this *SocketController) ToConn() {

	ws, err := upgrader.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil)

	if err != nil {
		fmt.Println(err)
	}

	if _, ok := clients[ws]; !ok {

		typing := this.GetString("type")
		hotel, err := this.GetInt("hotel")

		if err != nil {
			fmt.Print(err)
			this.StopRun()
		}

		if typing == "" || hotel == 0 {
			if err := ws.WriteJSON(response); err != nil {
				fmt.Print(err)
				ws.Close()
			}
		}

		clients[ws] = Info{typing, hotel, nil}

		if err := ws.WriteJSON(response); err != nil {
			fmt.Print(err)
			ws.Close()
			delete(clients, ws)
		}
	}

	_, data, err := ws.ReadMessage()
	if err != nil {
		if err := ws.WriteJSON(response); err != nil {
			fmt.Print(err)
			ws.Close()
			delete(clients, ws)
		}
	}

	byteData := []byte(data)

	var jsonData interface{}

	err = json.Unmarshal(byteData, &jsonData)

	if err != nil {
		response.Msg = "数据格式错误"
		if err := ws.WriteJSON(response); err != nil {
			fmt.Print(err)
			ws.Close()
			delete(clients, ws)
		}
	}

	fmt.Println(len(byteData))

	if info.typing == "" || info.hotel == 0 || len(info.results) == 0 {
		response.Msg = "数据格式错误"
		if err := ws.WriteJSON(response); err != nil {
			fmt.Print(err)
			ws.Close()
			delete(clients, ws)
		}
	}

	response.Result = info.results
	response.Msg = "success"

	for key, val := range clients {

		if key == ws || info.typing != val.typing || info.hotel != val.hotel {
			continue
		}
		if err := key.WriteJSON(response); err != nil {
			fmt.Print(err)
			key.Close()
			delete(clients, ws)
		}
	}
}

func (this *SocketController) ToSend() {

	typing := this.GetString("type")

	if typing == "" {
		response.Code = 422
		response.Msg = "类型不能为空"
		this.Data["json"] = response
		this.ServeJSON()
		this.StopRun()
	}

	hotel, err := this.GetInt("hotel")

	fmt.Println(hotel)

	if err != nil {
		response.Code = 422
		response.Msg = "酒店为空或格式错误"
		this.Data["json"] = response
		this.ServeJSON()
		this.StopRun()
	}

	var content = make(map[string]string)

	this.Ctx.Input.Bind(&content, "content")

	if len(content) < 1 {
		response.Code = 422
		response.Msg = "发送内容不能为空"
		this.Data["json"] = response
		this.ServeJSON()
		this.StopRun()
	}

	this.Data["json"] = response
	this.ServeJSON()

	response.Result = content

	for key, val := range clients {

		if val.typing != typing || val.hotel != hotel {
			continue
		}

		if err := key.WriteJSON(response); err != nil {
			fmt.Print(err)
			key.Close()
			delete(clients, key)
		}
	}
}
