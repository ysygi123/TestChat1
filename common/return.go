package common

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ReturnResponse(c *gin.Context, code int, statusCode int, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    statusCode,
		Message: message,
		Data:    data,
	})
}

func WebSocketReturn(conn *websocket.Conn, code int, message string, data interface{}) {
	j, _ := json.Marshal(Response{
		code,
		message,
		data,
	})
	_ = conn.WriteMessage(websocket.TextMessage, j)
}
