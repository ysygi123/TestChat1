package webcontroller

import (
	"TestChat1/common"
	"TestChat1/db/mysql"
	"TestChat1/db/redis"
	"TestChat1/vaildate/uservalidate"
	"github.com/gin-gonic/gin"
	"time"
)

func Login(c *gin.Context) {
	var userParams uservalidate.LoginValidate
	err := common.AutoValidate(c, &userParams)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	rows, err := mysql.DB.Query("select uid,passwd,username from user where username=? limit 1", userParams.Username)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	//将数据转为map
	userData, err := mysql.GetOneRow(rows)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	//验证密码
	if common.GetMD5Data(userParams.Password) != userData["passwd"] {
		common.ReturnResponse(c, 200, 400, "密码错误", nil)
		return
	}
	//获取session
	session := common.GetSession(userParams.Username)
	rec := redis.RedisPool.Get()
	//判断是否登录
	replay, err := rec.Do("GET", "uidlogin:"+userData["uid"])
	if err != nil {
		common.ReturnResponse(c, 200, 400, "取出缓存错误", nil)
		return
	}
	if replay != nil {
		common.ReturnResponse(c, 200, 400, "请勿重复登陆", nil)
		return
	}
	//设置基础信息
	_, err = rec.Do("HMSET", session, "uid", userData["uid"], "username", userData["username"])
	if err != nil {
		common.ReturnResponse(c, 200, 400, "设置token错误", nil)
		return
	}
	//设置是否登录
	_, err = rec.Do("SET", "uidlogin:"+userData["uid"], uint64(time.Now().Unix()))
	if err != nil {
		common.ReturnResponse(c, 200, 400, "设置登录位错误", nil)
		return
	}
	/*uid, _ := strconv.Atoi(userData["uid"])
	cli := websocket.NewClient(c.ClientIP(), uid, uint64(time.Now().Unix()), session)
	websocket.ClientMangerInstance.AddClient(uid, cli)*/
	common.ReturnResponse(c, 200, 200, "登陆成功", map[string]string{
		"session":  session,
		"uid":      userData["uid"],
		"username": userData["username"],
	})
	return
}
