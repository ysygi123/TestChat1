package websocket

import (
	"TestChat1/db/mysql"
	"TestChat1/db/redis"
	"TestChat1/router/websocketroute"
)

var ClientMangerInstance *ClientManger

func WebsocketInin() {
	ClientMangerInstance = new(ClientManger)
	mysql.NewMysqlDB()
	redis.NewRedisDB()
	ClientMangerInstance = NewClientManger()
	websocketroute.NewWebSocketRoute()
}

func WebSocketStart() {

}
