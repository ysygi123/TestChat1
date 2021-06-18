package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type response struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

func ReturnResponse(c *gin.Context, code int, statusCode int, message string, data interface{})  {
	c.JSON(http.StatusOK, response{
		Code:statusCode,
		Message:message,
		Data:data,
	})
}