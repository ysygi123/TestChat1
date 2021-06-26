package webcontroller

import (
	"TestChat1/common"
	"TestChat1/db/mysql"
	"TestChat1/model/message"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetMessageList(c *gin.Context) {
	uidStr := c.Param("uid")
	uid, _ := strconv.Atoi(uidStr)
	if uid == 0 {
		common.ReturnResponse(c, 200, 400, "id错误", nil)
		return
	}
	rows, err := mysql.DB.Query("SELECT `id`,`from_id`,`message_content`,`message_type`,`created_time`,`update_time`,`message_num`,`message_id` FROM `message_list` WHERE `uid`=? AND `is_del`=1", uid)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	var messageList []message.MessageList
	var tmp message.MessageList
	for rows.Next() {
		err := rows.Scan(&tmp.Id, &tmp.FromId, &tmp.MessageContent, &tmp.MessageType, &tmp.CreatedTime, &tmp.UpdateTime, &tmp.MessageNum, &tmp.MessageId)
		if err != nil {
			common.ReturnResponse(c, 200, 400, err.Error(), nil)
			return
		}
		messageList = append(messageList, tmp)
	}
	common.ReturnResponse(c, 200, 200, "success", messageList)
}
