package main

import (
	"TestChat1/router/webroute"
	"TestChat1/servers/web"
	"github.com/gorilla/websocket"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main()  {
	webroute.SetWebRoute()
	_ = web.GinEniger.Run("192.168.199.112:8088")
}

