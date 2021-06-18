package websocket

import "github.com/gorilla/websocket"

type client struct {
	Ip             string
	Id             int
	HeartBreath    uint64
	BelongServerID int
	WebSocketConn  websocket.Conn
}

func NewClient(ip string, id int, heartBreath uint64, websocketconn websocket.Conn) *client {
	return &client{
		Ip:             ip,
		Id:             id,
		HeartBreath:    heartBreath,
		BelongServerID: 1,
		WebSocketConn:  websocketconn,
	}
}
