package testTool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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
