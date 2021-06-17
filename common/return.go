package common

import "github.com/gin-gonic/gin"

type response struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

func ReturnResponse(c *gin.Context, code int, message string, data interface{})  {
	c.JSON(code, response{
		Code:code,
		Message:message,
		Data:data,
	})
}