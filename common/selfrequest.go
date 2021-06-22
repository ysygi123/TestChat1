package common

import "encoding/json"

type WebSocketRequest struct {
	Cmd  string                 `json:"cmd"`
	Body map[string]interface{} `json:"body"`
}

func GetNewWebSocketRequest(cmd string) *WebSocketRequest {
	return &WebSocketRequest{
		Cmd:  cmd,
		Body: make(map[string]interface{}),
	}
}

func GetJsonByteData(w *WebSocketRequest) []byte {
	j, _ := json.Marshal(w)
	return j
}
