package websocketcontroller

import (
	"TestChat1/common"
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
)

func FirstPage(w http.ResponseWriter, req *http.Request) {
	upgrade := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrade.Upgrade(w, req, nil)
	if err != nil {
		http.NotFound(w, req)
		return
	}
	returnData := common.Response{
		Code:    200,
		Message: "你成功了",
		Data: map[string]string{
			"IP": conn.RemoteAddr().String(),
		},
	}
	b, err := json.Marshal(returnData)
	err = conn.WriteMessage(0, b)
	if err != nil {
		http.NotFound(w, req)
		return
	}
}
