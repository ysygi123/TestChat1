package main

import (
	"TestChat1/controller/websocketcontroller"
	"TestChat1/db/mysql"
	"TestChat1/db/redis"
	"TestChat1/router/webroute"
	"TestChat1/router/websocketroute"
	"TestChat1/servers/web"
	"TestChat1/servers/websocket"
	"net/http"
)

func main() {
	mysql.NewMysqlDB()
	redis.NewRedisDB()
	websocket.ClientMangerInstanceInit()
	websocketroute.WebSocketRouteMangerInit()
	go func() {
		http.HandleFunc("/ws", websocketcontroller.FirstPage)
		http.ListenAndServe(":8087", nil)
	}()
	webroute.SetWebRoute()
	_ = web.GinEniger.Run(":8088")
}
