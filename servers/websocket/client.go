package websocket

import (
	"TestChat1/common"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"runtime/debug"
)

type Client struct {
	MessageChannel chan *common.WebSocketRequest
	Uid            int
	BelongServerID int
	HeartBreath    uint64
	Ip             string
	WebSocketConn  *websocket.Conn
	IsAuth         bool
}

func NewClient(ip string, uid int, heartBreath uint64, websocketConn *websocket.Conn) *Client {
	return &Client{
		MessageChannel: make(chan *common.WebSocketRequest, 1000),
		Ip:             ip,
		Uid:            uid,
		HeartBreath:    heartBreath,
		BelongServerID: 1,
		WebSocketConn:  websocketConn,
		IsAuth:         false,
	}
}

//读消息
func (this *Client) ReadData() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("我输错啦 这里是读", string(debug.Stack()), r)
		}
	}()

	defer func() {
		close(this.MessageChannel)
	}()

	for {
		mesType, mesg, err := this.WebSocketConn.ReadMessage()
		if err != nil {
			this.HandleErrorData(err, mesType, "reload")
			ClientMangerInstance.CloseChan <- this.Uid
			break
		}
		this.handleData(mesType, mesg)
	}
}

//向conn端写消息
func (this *Client) WriteData() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("我输错啦 这里是写写写", string(debug.Stack()), r)
		}
	}()
	//向manager的这个发送关闭信息
	defer func() {
		ClientMangerInstance.CloseChan <- this.Uid
	}()

	for {
		select {
		case msg, ok := <-this.MessageChannel:
			if !ok {
				//关闭连接
				return
			}

			_ = this.WebSocketConn.WriteMessage(websocket.TextMessage, common.GetJsonByteData(msg))
		}
	}
}

//发送消息用的 对外提供投送消息
func (this *Client) SendMsg(msg *common.WebSocketRequest) {

	if this == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("我已经没了", r, string(debug.Stack()))
		}
	}()
	this.MessageChannel <- msg
}

//客户端操作数据
func (this *Client) handleData(mesType int, mesg []byte) {
	s := new(common.WebSocketRequest)
	err := json.Unmarshal(mesg, s)
	if err != nil {
		this.HandleErrorData(err, mesType, "reload")
		return
	}
	hFunction, err := WebSocketRouteManger.GetHandler(s.Cmd)
	if err != nil {
		this.HandleErrorData(err, mesType, "reload")
		return
	}
	if this.IsAuth == false {
		err = errors.New("先认证")
		this.HandleErrorData(err, mesType, "reAuth")
		return
	}
	hFunction(this, &s.Body)
	//hFunction(this, []byte(strData.Message))
}

//错误操作
func (this *Client) HandleErrorData(err error, mesType int, cmd string) {
	s := common.GetNewWebSocketRequest(cmd)
	s.Body["err"] = err.Error()
	this.WebSocketConn.WriteMessage(mesType, common.GetJsonByteData(s))
	ClientMangerInstance.CloseChan <- this.Uid
}
