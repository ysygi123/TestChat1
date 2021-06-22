package common

type WebSocketRequest struct {
	Cmd     string            `json:"cmd"`
	Message map[string]string `json:"message"`
}
