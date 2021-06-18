package main

import (
	"TestChat1/router/webroute"
	"TestChat1/router/websocketroute"
	"TestChat1/servers/web"
)

func main() {
	go websocketroute.StartWebSocketRoute()
	webroute.SetWebRoute()
	_ = web.GinEniger.Run(":8088")
}
