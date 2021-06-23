package webcontroller

import (
	"TestChat1/common"
	"TestChat1/db/redis"
	"TestChat1/model/message"
	"TestChat1/vaildate/messagevalidate"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"
)

//发送消息
func SendMessage(c *gin.Context) {
	var messageParams messagevalidate.Message
	err := common.AutoValidate(c, &messageParams)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	pm := &message.PipelineMessage{
		MessageType: 1,
		MessageBody: message.Message{
			SendUid:        messageParams.SendUid,
			ReceiveUid:     messageParams.ReceiveUid,
			MessageContent: messageParams.MessageContent,
			CreatedTime:    uint64(time.Now().Unix()),
		},
	}
	jsonMessage, err := json.Marshal(pm)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	//发送消息到队列
	rec := redis.RedisPool.Get()
	_, err = rec.Do("LPUSH", "message_queue", jsonMessage)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	common.ReturnResponse(c, 200, 200, "发送成功", nil)
}
