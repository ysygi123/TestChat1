package websocket

import (
	"TestChat1/db/mysql"
	"TestChat1/db/redis"
)

var ClientMangerInstance *ClientManger

func WebsocketInin() {
	ClientMangerInstance = new(ClientManger)
	mysql.NewMysqlDB()
	redis.NewRedisDB()
	ClientMangerInstance = NewClientManger()
}

func WebSocketStart() {

}
