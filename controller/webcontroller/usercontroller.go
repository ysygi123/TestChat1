package webcontroller

import (
	"TestChat1/common"
	"TestChat1/db/mysql"
	"TestChat1/vaildate/uservalidate"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context)  {
	var userParams uservalidate.LoginValidate
	err := common.AutoValidate(c, &userParams)
	if err != nil{
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	rows, err := mysql.DB.Query("select uid,passwd from user where username=? limit 1", userParams.Username)
	if err != nil{
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	userData, err := mysql.GetOneRow(rows)
	if err != nil{
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	fmt.Println(common.GetMD5Data(userParams.Password), userData["passwd"])
	if common.GetMD5Data(userParams.Password) != userData["passwd"] {
		common.ReturnResponse(c, 200, 400, "密码错误", nil)
		return
	}
	session := common.GetSession(userParams.Username)
	common.ReturnResponse(c, 200, 200, "登陆成功", map[string]string{"session" : session})
	return
}