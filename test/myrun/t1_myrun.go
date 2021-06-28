package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"

	"io/ioutil"
	"net/http"
)

func main() {
	imitate(1)
}

func imitate(uid int) {
	//登录修改
	urlLogin := "http://127.0.0.1:8088/user/Login"
	requestBody := fmt.Sprintf(`{"username":"%d","password":"%d"}`, uid+9, uid+9)
	loginreturnmap := myselfPostRequest(urlLogin, requestBody, nil)
	datamapinterface := loginreturnmap["data"]
	datamap := datamapinterface.(map[string]interface{})
	session := datamap["session"].(string)
	//websocket
	websocketUrl := fmt.Sprintf("ws://127.0.0.1:8087/ws?uid=%d", uid)
	ws, err := websocket.Dial(websocketUrl, "", "http://127.0.0.0:8088")
	if err != nil {
		fmt.Println(err, session)
	}
	msg := make([]byte, 1024)
	_, err = ws.Read(msg)
	fmt.Println("这个是收到的消息哇", string(msg))
	//验证了
	authUrl := "http://127.0.0.1:8088/user/AuthClient"
	requestBody = fmt.Sprintf(`{"uid":%d,"session":"%s"}`, uid, session)
	myselfPostRequest(authUrl, requestBody, map[string]string{"session": session})

}

func myselfPostRequest(url, requestBody string, header map[string]string) map[string]interface{} {

	jsonStr := []byte(requestBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("request error", err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	if len(header) > 0 {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("client get resp", err)
		return nil
	}

	body, _ := ioutil.ReadAll(resp.Body)
	returnMap := make(map[string]interface{})
	json.Unmarshal(body, &returnMap)
	return returnMap
}
