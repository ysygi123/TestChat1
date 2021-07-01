package websocket

import (
	"TestChat1/common"
	"github.com/gorilla/websocket"
	"time"
)

//私法消息
func SendMessageToOneUser(client *Client, message *map[string]interface{}) {
	err := common.CheckWebSocketParamsIsUnEmpty([]string{"uid", "message"}, message)
	if err != nil {
		client.HandleErrorData(err, websocket.TextMessage, "resend")
	}
	uid := (*message)["uid"].(int)
	otherC, err := ClientMangerInstance.GetClient(uid)
	if err != nil {
		if otherC == nil {
			return
		}
		client.HandleErrorData(err, websocket.TextMessage, "reload")
	}

	wsq := common.GetNewWebSocketRequest("GetMessage")
	wsq.Body["message"] = (*message)["message"]

	otherC.WebSocketConn.WriteMessage(websocket.TextMessage, common.GetJsonByteData(wsq))
}

func HeartBreath(client *Client, message *map[string]interface{}) {
	t := uint64(time.Now().Unix())
	ClientMangerInstance.SetClientHeartBreath(client.Uid, t)
}
