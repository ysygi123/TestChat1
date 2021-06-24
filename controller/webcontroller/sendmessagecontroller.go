package webcontroller

import (
	"TestChat1/common"
	"TestChat1/model/message"
	"TestChat1/servers/factory/messageDispatch"
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"time"
)

//发送消息
func SendMessage(c *gin.Context) {
	msg := &message.Message{}
	err := c.ShouldBindJSON(msg)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	msg.CreatedTime = uint64(time.Now().Unix())

	fmt.Printf("奇奇怪怪 %+v\n\n", msg)

	mfc, err := messageDispatch.CreateMessage(map[string]interface{}{
		"messageType": msg.MessageType,
	})

	fmt.Println(reflect.TypeOf(mfc))

	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	err = mfc.PushMessage(msg)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	common.ReturnResponse(c, 200, 200, "发送成功", nil)
}
