package main

import (
	"TestChat1/controller/websocketcontroller"
	"TestChat1/db/mysql"
	"TestChat1/db/redis"
	"TestChat1/router/webroute"
	"TestChat1/servers/backtask"
	"TestChat1/servers/web"
	"TestChat1/servers/websocket"
	"net/http"
)

func main() {
	mysql.NewMysqlDB()
	redis.NewRedisDB()
	websocket.ClientMangerInstanceInit()
	websocket.WebSocketRouteMangerInit()
	websocket.WebSocketRouteManger.AllRegisterRoute()
	go websocket.ClientMangerInstance.LoopToKillChild()
	go backtask.TaskConsumeMessage()
	go func() {
		http.HandleFunc("/ws", websocketcontroller.FirstPage)
		http.ListenAndServe(":8087", nil)
	}()
	webroute.SetWebRoute()
	_ = web.GinEniger.Run(":8088")
}
