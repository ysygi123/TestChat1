package messageChild

import (
	"TestChat1/common"
	"TestChat1/db/mysql"
	"TestChat1/model/message"
	"TestChat1/servers/websocket"
	"fmt"
	"strconv"
	"time"
)

type UserMessage struct {
}

func (this *UserMessage) AddMessage(messageData *message.PipelineMessage) error {
	m := messageData.MessageBody.(map[string]interface{})
	msg := message.Message{
		CreatedTime:    uint64(m["created_time"].(float64)),
		ReceiveUid:     int(m["receive_uid"].(float64)),
		SendUid:        int(m["send_uid"].(float64)),
		MessageContent: m["message_content"].(string),
	}
	//即时发送消息
	go func() {
		c, err := websocket.ClientMangerInstance.GetClient(msg.ReceiveUid)
		if err != nil {
			fmt.Println("错误")
			return
		}
		wmsg := common.GetNewWebSocketRequest("GetMessage")
		wmsg.Body = map[string]interface{}{
			"message_content": msg.MessageContent,
		}
		_ = c.WebSocketConn.WriteMessage(1, common.GetJsonByteData(wmsg))
	}()

	//处理消息入库
	res, err := mysql.DB.Exec("BEGIN")
	msg.CreatedTime = uint64(time.Now().Unix())
	res, err = mysql.DB.Exec(
		"INSERT INTO `message`"+
			"(`message_content`,`send_uid`,`receive_uid`,`created_time`) "+
			"VALUES (?,?,?,?)", msg.MessageContent, msg.SendUid, msg.ReceiveUid, msg.CreatedTime)
	if err != nil {
		fmt.Println("clientManager line 126: ", res, err)
		mysql.DB.Exec("ROLLBACK")
	}

	resq, err := mysql.DB.Query("SELECT `id`,`message_num`,`uid`,`from_id` FROM `message_list` WHERE (`from_id`=? OR `uid`=?) AND `message_type`=1 AND `is_del`=1 limit 1", msg.SendUid, msg.SendUid)

	if err != nil {
		fmt.Println("clientManager line 131", resq, err)
	}

	tmpMsg := message.MessageList{}

	iSendFlag := false
	iReceiveFlag := false
	fmt.Println(iSendFlag, iReceiveFlag)
	for resq.Next() {
		err = resq.Scan(&tmpMsg.Id, &tmpMsg.MessageNum, &tmpMsg.Uid, &tmpMsg.FromId)
		if err != nil {
			fmt.Println("clientManager line 61: ", err)
			return nil
		}
		//那么这条是我发的
		if tmpMsg.FromId == msg.SendUid {
			//将标变为正确
			iSendFlag = true

		}
	}
	resqMap, err := mysql.GetOneRow(resq)
	//设置可插入list表的消息
	msgTitle := []rune(msg.MessageContent)
	lmsgTitle := len(msgTitle)
	msgcontent := ""
	if lmsgTitle >= 50 {
		msgcontent = string(msgTitle[0:47]) + "..."
	} else {
		msgcontent = msg.MessageContent
	}

	if idStr, ok := resqMap["id"]; ok == false { //新增
		res, err := mysql.DB.Exec("INSERT INTO `message_list`"+
			"(`uid`,`from_id`,`message_content`,`message_type`,`created_time`,`update_time`,`message_num`)"+
			"VALUES (?,?,?,?,?,?,?)", msg.ReceiveUid, msg.SendUid, msgcontent, 1, msg.CreatedTime, msg.CreatedTime, 1)
		if err != nil {
			fmt.Println("clientManager line 150", res, err)
			mysql.DB.Exec("ROLLBACK")
		}
	} else { //修改
		id, _ := strconv.Atoi(idStr)
		res, err := mysql.DB.Exec("UPDATE `message_list` SET "+
			"`message_content`=?,update_time=?,`message_num`=`message_num`+1 "+
			"WHERE id=?",
			msgcontent, msg.CreatedTime, id)
		if err != nil {
			fmt.Println("clientManager line 159", res, err)
			mysql.DB.Exec("ROLLBACK")
		}
	}
	res, err = mysql.DB.Exec("COMMIT")

	return nil
}
