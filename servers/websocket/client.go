package websocket

import (
	"TestChat1/common"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

type Client struct {
	ErrChan        chan struct{}
	Uid            int
	BelongServerID int
	HeartBreath    uint64
	Ip             string
	WebSocketConn  *websocket.Conn
}

func NewClient(ip string, uid int, heartBreath uint64, websocketConn *websocket.Conn) *Client {
	return &Client{
		ErrChan:        make(chan struct{}, 1),
		Ip:             ip,
		Uid:            uid,
		HeartBreath:    heartBreath,
		BelongServerID: 1,
		WebSocketConn:  websocketConn,
	}
}

func (this *Client) ReadData() {
	defer func() { this.ErrChan <- struct{}{} }()
	for {
		mesType, mesg, err := this.WebSocketConn.ReadMessage()
		if err != nil {
			fmt.Println("这里是协程测试出现错误的情况", mesType, mesg, err, string(mesg))
			break
		}
		fmt.Println("这里是协程测试", mesType, mesg, err, string(mesg))
		this.handleData(mesType, mesg)
	}
}

func (this *Client) handleData(mesType int, mesg []byte) {
	strData := &common.WebSocketRequest{}
	err := json.Unmarshal(mesg, strData)
	fmt.Println(strData)
	if err != nil {
		this.WebSocketConn.WriteMessage(mesType, []byte("出现错误1 : "+err.Error()))
	}

	//hFunction(this, []byte(strData.Message))
}

func (this *Client) WriteData() {

}
