package main

import (
	"TestChat1/test/myrun/testTool"
	"fmt"
	"golang.org/x/net/websocket"
	"sync"
	"time"
)

func main() {
	Register(1)
	/*fatherWG := new(sync.WaitGroup)
	for j := 1; j < 10; j++ {
		fatherWG.Add(1)
		go imitate(j, fatherWG)
	}
	fatherWG.Wait()*/
	//AddGroup()
}

var prefixHttpUrl = "http://192.168.3.36:8088"
var prefixWsUrl = "ws://192.168.3.36:8087"

func imitate(uid int, fatherWG *sync.WaitGroup) {
	//登录修改
	defer fatherWG.Done()
	urlLogin := prefixHttpUrl + "/user/Login"
	requestBody := fmt.Sprintf(`{"username":"%d","password":"%d"}`, uid, uid)
	loginreturnmap := testTool.MyselfPostRequest(urlLogin, requestBody, nil)
	if loginreturnmap == nil {
		fmt.Println("uid", uid, "报错了 没有登录返回map")
		return
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
	//验证
	sendStr := fmt.Sprintf(`{"cmd":"Auth","body":{"session":"%s"}}`, session)
	n, err := ws.Write([]byte(sendStr))
	fmt.Println("uid : ", uid, "查看", n, err)
	_, err = ws.Read(msg)
	fmt.Println("uid : ", uid, "查看验证完的数据是什么", string(msg), "看看发送的数据是什么", sendStr)
	wg := new(sync.WaitGroup)
	//发心跳
	go func() {
		for {
			_, _ = ws.Write([]byte(fmt.Sprintf(`{"cmd":"HeartBreath"}`)))
			time.Sleep(time.Second * 60)
		}
	}()
	TestHttpTool := new(testTool.TestHttp)
	groupChatUrl := prefixHttpUrl + "/message/SendMessage"
	requestBody = fmt.Sprintf(`{"send_uid":%d,"chat_id":1,"message_content":"你好我叫%d当前时间是%d","message_type":2}`,
		uid, uid, time.Now().UnixNano()/1e6)
	TestHttpTool.NewTestHttp(groupChatUrl, requestBody, map[string]string{"session": session})

	go func() {
		for num := 0; num < 2; num++ {
			wg.Add(1)
			go func() {
				res := TestHttpTool.SendRequest()
				fmt.Println("uid : ", uid, "查看发送聊天消息的结果是什么 : ", res)
				wg.Done()
			}()
		}
	}()
	i := 1
	for {
		_, err := ws.Read(msg)
		i++
		if err != nil {
			fmt.Println("uid", uid, "这也是个奇奇怪怪的error", err)
			break
		}
		fmt.Println("我是uid : ", uid, "我收到的消息是 : ", string(msg), i)
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

func Register(i int) {
	for k := i; k < 501; k++ {
		jsonStr := fmt.Sprintf(`{"username":"%d","passwd":"%d"}`, k, k)
		fmt.Println(jsonStr)
		res := testTool.MyselfPostRequest("http://127.0.0.1:8088/user/Register", jsonStr, nil)
		fmt.Println(res)
	}
}
