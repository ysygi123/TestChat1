package web

import "github.com/gin-gonic/gin"

var GinEniger *gin.Engine

func init()  {
	GinEniger = gin.New()
}

