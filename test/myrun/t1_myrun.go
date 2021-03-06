package main

import (
	"TestChat1/test/myrun/testTool"
	"fmt"
	"golang.org/x/net/websocket"
	"sync"
	"time"
)

func main() {
	//Register(1)
	fatherWG := new(sync.WaitGroup)
	c := testTool.GetNewChan(499)
	for j := 2; j < 10001; j++ {
		fatherWG.Add(1)
		go imitate(j, fatherWG, c)
	}
	fatherWG.Wait()
	//AddGroup()
}

var msgIdTime = uint64(time.Now().Unix())
var prefixHttpUrl = "http://192.168.199.112:8088"
var prefixWsUrl = "ws://192.168.199.112:8087"

func imitate(uid int, fatherWG *sync.WaitGroup, cLock *testTool.SendChan) {
	//登录修改
	defer fatherWG.Done()
	cLock.Lock()
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
	cLock.Unlock()
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
	/*TestHttpTool := new(testTool.TestHttp)
	groupChatUrl := prefixHttpUrl + "/message/SendMessage"
	requestBody = fmt.Sprintf(`{"send_uid":%d,"chat_id":1,"message_content":"你好我叫%d当前时间是%d","message_type":2}`,
		uid, uid, msgIdTime)
	TestHttpTool.NewTestHttp(groupChatUrl, requestBody, map[string]string{"session": session})

	go func() {
		for num := 0; num < 3; num++ {
			wg.Add(1)
			go func() {
				res := TestHttpTool.SendRequest()
				fmt.Println("uid : ", uid, "查看发送聊天消息的结果是什么 : ", res)
				wg.Done()
			}()
		}
	}()*/
	i := 1
	for {
		n, err := ws.Read(msg)
		i++
		if err != nil {
			fmt.Println("uid", uid, "这也是个奇奇怪怪的error", err)
			break
		}
		fmt.Println("我是uid : ", uid, "我收到的消息是 : ", string(msg[:n]), i)
	}
	wg.Wait()
}

func AddGroup() {
	url := "http://127.0.0.1:8088/group/AddToGroupCommit"
	c := testTool.GetNewChan(499)
	wg := new(sync.WaitGroup)
	for i := 1; i < 10001; i++ {
		jsonStr := fmt.Sprintf(`{"uid":%d,"group_id":1}`, i)
		wg.Add(1)
		c.Lock()
		go func() {
			AddGroupCommitMap := testTool.MyselfPostRequest(url, jsonStr, nil)
			fmt.Println(AddGroupCommitMap)
			c.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
}

func Register(i int) {
	for k := 501; k < 10001; k++ {
		jsonStr := fmt.Sprintf(`{"username":"%d","passwd":"%d"}`, k, k)
		fmt.Println(jsonStr)
		res := testTool.MyselfPostRequest("http://127.0.0.1:8088/user/Register", jsonStr, nil)
		fmt.Println(res)
	}
}
