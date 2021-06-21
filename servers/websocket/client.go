package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type client struct {
	ErrChan        chan struct{}
	Uid            int
	BelongServerID int
	HeartBreath    uint64
	Ip             string
	WebSocketConn  *websocket.Conn
}

func NewClient(ip string, uid int, heartBreath uint64, websocketConn *websocket.Conn) *client {
	return &client{
		ErrChan:        make(chan struct{}, 1),
		Ip:             ip,
		Uid:            uid,
		HeartBreath:    heartBreath,
		BelongServerID: 1,
		WebSocketConn:  websocketConn,
	}
}

func (this *client) ReadData() {
	defer func() { this.ErrChan <- struct{}{} }()
	for {
		mesType, mesg, err := this.WebSocketConn.ReadMessage()
		if err != nil {
			fmt.Println("这里是协程测试出现错误的情况", mesType, mesg, err, string(mesg))
			break
		}
		fmt.Println("这里是协程测试", mesType, mesg, err, string(mesg))
	}
}

func (this *client) WriteData() {

}
