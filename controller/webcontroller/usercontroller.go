package webcontroller

import (
	"TestChat1/common"
	"TestChat1/vaildate/uservalidate"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context)  {
	var uservalidatee uservalidate.LoginValidate
	if err := c.ShouldBindJSON(&uservalidatee); err != nil{
		fmt.Printf("%v")
		common.ReturnResponse(c, 400, err.Error(), nil)
		return
	}

}