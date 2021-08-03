package testTool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TestHttp struct {
	Req *http.Request
	Cli *http.Client
}

func (this *TestHttp) NewTestHttp(url, requestBody string, header map[string]string) {
	jsonStr := []byte(requestBody)
	var err error
	this.Req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	if err != nil {
		fmt.Println("这个是错误啊啊啊啊", err)
		return
	}
	this.Req.Header.Set("Content-Type", "application/json")
	this.Req.Header.Set("Connection", "keep-alive")

	if len(header) > 0 {
		for k, v := range header {
			this.Req.Header.Set(k, v)
		}
	}
	this.Cli = &http.Client{}
}

func (this *TestHttp) SendRequest() map[string]interface{} {
	resp, err := this.Cli.Do(this.Req)
	if err != nil {
		fmt.Println("client get resp", err)
		return nil
	}
	body, _ := ioutil.ReadAll(resp.Body)
	returnMap := make(map[string]interface{})
	json.Unmarshal(body, &returnMap)
	return returnMap
}

//----------------------------------
func MyselfPostRequest(url, requestBody string, header map[string]string) map[string]interface{} {

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

type SendChan struct {
	ch chan struct{}
}

func GetNewChan(numOfChan int) *SendChan {
	s := &SendChan{}
	s.ch = make(chan struct{}, numOfChan)
	return s
}

func (this *SendChan) Lock() {
	this.ch <- struct{}{}
}

func (this *SendChan) Unlock() {
	<-this.ch
}
