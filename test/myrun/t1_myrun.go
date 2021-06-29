package main

import (
	"TestChat1/test/myrun/testTool"
	"fmt"
	"golang.org/x/net/websocket"
	"sync"
	"time"
)

func main() {
	/*fatherWG := new(sync.WaitGroup)
	for j := 1; j < 3; j++ {
		fatherWG.Add(1)
		go imitate(j, fatherWG)
	}
	fatherWG.Wait()*/
	AddGroup()
}

var prefixHttpUrl = "http://127.0.0.1:8088"
var prefixWsUrl = "ws://127.0.0.1:8087"

func imitate(uid int, fatherWG *sync.WaitGroup) {
	//登录修改
	defer fatherWG.Done()
	urlLogin := prefixHttpUrl + "/user/Login"
	requestBody := fmt.Sprintf(`{"username":"%d","password":"%d"}`, uid+9, uid+9)
	loginreturnmap := testTool.MyselfPostRequest(urlLogin, requestBody, nil)
	if loginreturnmap == nil {
		fmt.Println("uid", uid, "报错了 没有登录返回map")
	}
	datamapinterface := loginreturnmap["data"]
	if datamapinterface == nil {
		fmt.Println("uid", uid, "报错了 没有登录返回map2", loginreturnmap)
		return
	}
	datamap := datamapinterface.(map[string]interface{})
	session := datamap["session"].(string)
	//websocket
	websocketUrl := fmt.Sprintf(prefixWsUrl+"/ws?uid=%d", uid)
	ws, err := websocket.Dial(websocketUrl, "", "http://127.0.0.0:8088")
	if err != nil {
		fmt.Println(err, session)
		return
	}
	msg := make([]byte, 1024)
	_, err = ws.Read(msg)
	fmt.Println("这个是收到的消息哇", string(msg))
	//验证了
	authUrl := prefixHttpUrl + "/user/AuthClient"
	requestBody = fmt.Sprintf(`{"uid":%d,"session":"%s"}`, uid, session)
	authReturnMap := testTool.MyselfPostRequest(authUrl, requestBody, map[string]string{"session": session})
	authReturnCode := int(authReturnMap["code"].(float64))
	if authReturnCode != 200 {
		fmt.Println("错了验证")
		return
	}

	wg := new(sync.WaitGroup)
	for num := 0; num < 100; num++ {
		wg.Add(1)
		go func() {
			groupChatUrl := prefixHttpUrl + "/message/SendMessage"
			requestBody = fmt.Sprintf(`{"send_uid":%d,"group_id":1,"message_content":"你好我叫%d当前时间是%s","message_type":2}`,
				uid, uid, time.Now().Format("2006-01-02 15:04:05"))
			res := testTool.MyselfPostRequest(groupChatUrl, requestBody, map[string]string{"session": session})
			fmt.Println("查看发送聊天消息的结果是什么 : ", res)
			wg.Done()
		}()
	}
	wg.Wait()
}

func AddGroup() {
	url := "http://127.0.0.1:8088/group/AddToGroupCommit"
	wg := new(sync.WaitGroup)
	for i := 1; i < 501; i++ {
		jsonStr := fmt.Sprintf(`{"uid":%d,"group_id":1}`, i)
		wg.Add(1)
		go func() {
			AddGroupCommitMap := testTool.MyselfPostRequest(url, jsonStr, nil)
			fmt.Println(AddGroupCommitMap)
			wg.Done()
		}()
	}
	wg.Wait()
}
