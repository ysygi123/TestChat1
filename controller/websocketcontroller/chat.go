package websocketcontroller

import (
	"TestChat1/common"
	ww "TestChat1/servers/websocket"
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
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
	err = req.ParseForm()
	if err != nil {
		http.NotFound(w, req)
		return
	}
	returnData := common.Response{
		Code:    200,
		Message: "你成功了",
		Data: map[string]string{
			"cmd": "SendData",
		},
	}
	c := ww.NewClient(conn.RemoteAddr().String(), 1, uint64(time.Now().Unix()), conn)
	ww.ClientMangerInstance.AddClient(1, c)
	go c.ReadData()
	b, err := json.Marshal(returnData)
	err = conn.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		http.NotFound(w, req)
		return
	}
}
