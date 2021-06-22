package main

import (
	"TestChat1/router/webroute"
	"TestChat1/router/websocketroute"
	"TestChat1/servers/web"
	"TestChat1/servers/websocket"
)

func main() {
	websocket.WebsocketInin()
	websocketroute.NewWebSocketRoute()
	go websocketroute.StartWebSocketRoute()
	webroute.SetWebRoute()
	_ = web.GinEniger.Run(":8088")
}
