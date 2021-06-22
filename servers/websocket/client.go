package websocket

import (
	"TestChat1/common"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
)

type Client struct {
	ErrChan        chan struct{}
	Uid            int
	BelongServerID int
	HeartBreath    uint64
	Ip             string
	WebSocketConn  *websocket.Conn
	IsAuth         bool
}

func NewClient(ip string, uid int, heartBreath uint64, websocketConn *websocket.Conn) *Client {
	return &Client{
		ErrChan:        make(chan struct{}, 1),
		Ip:             ip,
		Uid:            uid,
		HeartBreath:    heartBreath,
		BelongServerID: 1,
		WebSocketConn:  websocketConn,
		IsAuth:         false,
	}
}

func (this *Client) ReadData() {
	defer func() { this.ErrChan <- struct{}{} }()
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

func (this *Client) WriteData() {

}
