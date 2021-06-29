package middleware

import (
	"TestChat1/common"
	"TestChat1/db/redis"
	"github.com/gin-gonic/gin"
)

func AuthSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := c.GetHeader("session")
		rec := redis.RedisPool.Get()
		defer rec.Close()
		res, err := rec.Do("HGET", session, "uid")
		if err != nil {
			common.ReturnResponse(c, 200, 400, err.Error(), nil)
			c.Abort()
			return
		}
		if res == nil {
			common.ReturnResponse(c, 200, 373, "认证失败", nil)
			c.Abort()
			return
		}
		//c.Next()
		return
	}
}
