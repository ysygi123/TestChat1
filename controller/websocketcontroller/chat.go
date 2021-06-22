package websocketcontroller

import (
	"TestChat1/common"
	"TestChat1/db/redis"
	ww "TestChat1/servers/websocket"
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
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
	query := req.URL.Query()

	//获取uid
	uidSlice, ok := query["uid"]
	if ok != true || len(uidSlice) != 1 {
		http.NotFound(w, req)
		return
	}
	uid, _ := strconv.Atoi(uidSlice[0])

	hasUid, err := CheckHasThisUid(uid)
	returnData := common.Response{
		Code:    200,
		Message: "你成功了",
		Data: map[string]string{
			"cmd": "SendData",
		},
	}
	//这个uid没有登录就要返回错误
	if hasUid == false {
		returnData.Code = 400
		returnData.Message = "这个uid没有登录啊 这个是redis查的 去你吗的"
		returnData.Data = map[string]string{
			"cmd": "reLogin",
		}
		b, _ := json.Marshal(returnData)
		err = conn.WriteMessage(websocket.TextMessage, b)
		_ = conn.Close()
		return
	}

	c := ww.NewClient(conn.RemoteAddr().String(), uid, uint64(time.Now().Unix()), conn)
	ww.ClientMangerInstance.AddClient(uid, c)
	go c.ReadData()
	b, err := json.Marshal(returnData)
	err = conn.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		http.NotFound(w, req)
		return
	}
}

func CheckHasThisUid(uid int) (bool, error) {
	rec := redis.RedisPool.Get()
	replay, err := rec.Do("GET", "uidlogin:"+strconv.Itoa(uid))
	if err != nil {
		return false, err
	}
	if replay == nil {
		return false, nil
	}
	return true, nil
}
