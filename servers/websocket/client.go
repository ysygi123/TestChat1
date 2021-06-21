package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type client struct {
	Ip             string
	Uid            int
	HeartBreath    uint64
	BelongServerID int
	WebSocketConn  websocket.Conn
}

func NewClient(ip string, uid int, heartBreath uint64, websocketconn websocket.Conn) *client {
	return &client{
		Ip:             ip,
		Uid:            uid,
		HeartBreath:    heartBreath,
		BelongServerID: 1,
		WebSocketConn:  websocketconn,
	}
}

func (this *client) ReadData() {
	for {
		mesType, mesg, err := this.WebSocketConn.ReadMessage()
		fmt.Println(mesType, mesg, err)
	}
}
