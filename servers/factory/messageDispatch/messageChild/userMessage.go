package messageChild

import (
	"TestChat1/common"
	"TestChat1/db/mysql"
	"TestChat1/model/message"
	"TestChat1/servers/websocket"
	"database/sql"
	"fmt"
	"time"
)

type UserMessage struct {
}

func (this *UserMessage) CommonHandle(ml *message.MessageList, resq *sql.Row, msg *message.Message, msgcontent string, isSelf bool) error {
	resq.Scan(&ml.Id, &ml.MessageNum, &ml.Uid, &ml.FromId)
	if ml.Id == 0 {
		if err := this.InsertData(msg, msgcontent); err != nil {
			return err
		}
	} else {
		if err := this.UpdateData(msg, msgcontent, ml.Id, isSelf); err != nil {
			return err
		}
	}
	return nil
}

//查询 我发出去的消息 也就是接收人的消息面板是不是要有红点
func (this *UserMessage) ISendMessage(ml *message.MessageList, msg *message.Message, msgcontent string) error {
	resq := mysql.DB.QueryRow("SELECT `id`,`message_num`,`uid`,`from_id` FROM `message_list` WHERE `from_id`=? AND `message_type`=1 limit 1", msg.SendUid)
	return this.CommonHandle(ml, resq, msg, msgcontent, true)
}

//查询 我收的消息里面是不是要有红点
func (this *UserMessage) IReceiverSendMessage(ml *message.MessageList, msg *message.Message, msgcontent string) error {
	resq := mysql.DB.QueryRow("SELECT `id`,`message_num`,`uid`,`from_id` FROM `message_list` WHERE `uid`=? AND `message_type`=1 limit 1", msg.SendUid)
	msg.SendUid, msg.ReceiveUid = msg.ReceiveUid, msg.SendUid
	return this.CommonHandle(ml, resq, msg, msgcontent, false)
}

//新增
func (this *UserMessage) InsertData(msg *message.Message, msgcontent string) error {
	res, err := mysql.DB.Exec("INSERT INTO `message_list`"+
		"(`uid`,`from_id`,`message_content`,`message_type`,`created_time`,`update_time`,`message_num`)"+
		"VALUES (?,?,?,?,?,?,?)", msg.ReceiveUid, msg.SendUid, msgcontent, 1, msg.CreatedTime, msg.CreatedTime, 1)
	if err != nil {
		fmt.Println("clientManager line 41", res, err)
		return err
	}
	return nil
}

//修改
func (this *UserMessage) UpdateData(msg *message.Message, msgcontent string, id int, isSelf bool) error {
	sqlsql := "UPDATE `message_list` SET " +
		"`message_content`=?,update_time=?,"
	if isSelf {
		sqlsql += "`message_num`=`message_num`+1,"
	}
	sqlsql += "`is_del`=1 WHERE id=?"
	res, err := mysql.DB.Exec(sqlsql, msgcontent, msg.CreatedTime, id)
	if err != nil {
		fmt.Println("clientManager line 54", res, err)
		return err
	}
	return nil
}

//主要逻辑
func (this *UserMessage) AddMessage(messageData *message.PipelineMessage) error {
	m := messageData.MessageBody.(map[string]interface{})
	msg := &message.Message{
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

	//设置可插入list表的消息
	msgTitle := []rune(msg.MessageContent)
	lmsgTitle := len(msgTitle)
	msgcontent := ""
	if lmsgTitle >= 50 {
		msgcontent = string(msgTitle[0:47]) + "..."
	} else {
		msgcontent = msg.MessageContent
	}
	ml := new(message.MessageList)
	if err = this.ISendMessage(ml, msg, msgcontent); err != nil {
		fmt.Println("ClientManager 116 : ", err)
		mysql.DB.Exec("ROLLBACK")
		return err
	}
	if err = this.IReceiverSendMessage(ml, msg, msgcontent); err != nil {
		fmt.Println("ClientManager 120 : ", err)
		mysql.DB.Exec("ROLLBACK")
		return err
	}
	res, err = mysql.DB.Exec("COMMIT")
	return nil
}
