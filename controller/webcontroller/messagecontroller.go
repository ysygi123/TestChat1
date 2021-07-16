package webcontroller

import (
	"TestChat1/common"
	"TestChat1/db/mysql"
	"TestChat1/model/message"
	"TestChat1/vaildate/messagevalidate"
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
	rows, err := mysql.DB.Query("SELECT `id`,`from_id`,`message_content`,`message_type`,`created_time`,`update_time`,`message_num` FROM `message_list` WHERE `uid`=? AND `is_del`=1", uid)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	var messageList []message.MessageList
	var tmp message.MessageList
	for rows.Next() {
		err := rows.Scan(&tmp.Id, &tmp.FromId, &tmp.MessageContent, &tmp.MessageType, &tmp.CreatedTime, &tmp.UpdateTime, &tmp.MessageNum)
		if err != nil {
			common.ReturnResponse(c, 200, 400, err.Error(), nil)
			return
		}
		messageList = append(messageList, tmp)
	}
	common.ReturnResponse(c, 200, 200, "success", messageList)
}

//获取普通私聊
func GetSelfChat(c *gin.Context) {
	gscv := &messagevalidate.GetSelfChatValidate{}
	err := common.AutoValidate(c, gscv)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	rows, err := mysql.DB.Query("select * from message where message_type=1 and "+
		"((receive_uid=? and send_uid=?) or "+
		"(receive_uid=? and send_uid=?)) created_time between ? and ?",
		gscv.ReceiveUid, gscv.SendUid, gscv.SendUid, gscv.ReceiveUid, gscv.StartTime, gscv.EndTime)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	returnMessage, err := mysql.GetManyRows(rows)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	common.ReturnResponse(c, 200, 200, "success", returnMessage)
}

//群消息
func GetGroupChat(c *gin.Context) {
	ggcv := &messagevalidate.GetGroupChatValidate{}
	err := common.AutoValidate(c, ggcv)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}

	rows, err := mysql.DB.Query("select * from message where message_type=2 and "+
		"group_id=? created_time between ? and ?", ggcv.GroupId, ggcv.StartTime, ggcv.EndTime)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	returnMessage, err := mysql.GetManyRows(rows)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	common.ReturnResponse(c, 200, 200, "success", returnMessage)
}
