package websocketcontroller

import (
	"TestChat1/common"
	"TestChat1/db/redis"
	ww "TestChat1/servers/websocket"
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
	wb := common.GetNewWebSocketRequest("")
	hasUid, err := CheckHasThisUid(uid)
	//这个uid没有登录就要返回错误
	if hasUid == false {
		wb.Cmd = "reLogin"
		err = conn.WriteJSON(wb)
		_ = conn.Close()
		return
	}
	//设置一个client
	c := ww.NewClient(conn.RemoteAddr().String(), uid, uint64(time.Now().Unix()), conn)
	ww.ClientMangerInstance.AddClient(uid, c)
	wb.Cmd = "SendAuth"
	conn.WriteJSON(wb)
	//阻塞读取json
	//设置超时等待
	//conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	err = conn.ReadJSON(wb)
	if err != nil {
		wb.Cmd = "reload"
		if conn != nil {
			conn.WriteJSON(wb)
			conn.Close()
		}
		return
	}
	//验证session
	if wb.Cmd != "Auth" {
		wb.Cmd = "reAuth"
		wb.Body["msg"] = "好像不是auth"
		conn.WriteJSON(wb)
		conn.Close()
		ww.ClientMangerInstance.CloseChan <- uid
		return
	}
	sessionInterface, ok := wb.Body["session"]
	if ok == false {
		wb.Cmd = "ErrorSession"
		wb.Body = map[string]interface{}{}
		conn.WriteJSON(wb)
		ww.ClientMangerInstance.CloseChan <- uid
		return
	}
	session := sessionInterface.(string)
	hasUid, err = checkSession(session, uid)

	if err != nil {
		wb.Cmd = "reAuth"
		wb.Body["msg"] = "这个是checksession之后错了" + err.Error()
		conn.WriteJSON(wb)
		ww.ClientMangerInstance.CloseChan <- uid
		return
	}

	wb.Cmd = "ok"
	wb.Body = nil
	conn.WriteJSON(wb)
	//开两个协程分别用来读取这个client发来的操作和给这个client发去指令
	go c.ReadData()
	go c.WriteData()
}

func CheckHasThisUid(uid int) (bool, error) {
	rec := redis.RedisPool.Get()
	replay, err := rec.Do("GET", "uidlogin:"+strconv.Itoa(uid))
	if err != nil {
		rec.Close()
		return false, err
	}
	if replay == nil {
		rec.Close()
		return false, nil
	}
	rec.Close()
	return true, nil
}

func checkSession(session string, uid int) (bool, error) {
	rec := redis.RedisPool.Get()
	replay, err := rec.Do("HGET", session, "uid")
	if err != nil {
		return false, err
	}
	if replay == nil {
		return false, nil
	}
	redisGetUidString := string([]byte(replay.([]uint8)))
	redisGetUid, err := strconv.Atoi(redisGetUidString)
	if err != nil {
		return false, nil
	}
	if redisGetUid != uid {
		return false, nil
	}
	err = ww.ClientMangerInstance.SetAuth(uid)
	if err != nil {
		return false, err
	}
	return true, nil
}
