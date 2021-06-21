package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type client struct {
	Uid            int
	BelongServerID int
	HeartBreath    uint64
	Ip             string
	WebSocketConn  websocket.Conn
}

func NewClient(ip string, uid int, heartBreath uint64) *client {
	return &client{
		Ip:             ip,
		Uid:            uid,
		HeartBreath:    heartBreath,
		BelongServerID: 1,
	}
}

func (this *client) ReadData() {
	for {
		mesType, mesg, err := this.WebSocketConn.ReadMessage()
		fmt.Println(mesType, mesg, err)
	}
}
