package messageChild

import (
	"TestChat1/common"
	"TestChat1/db/mysql"
	"TestChat1/db/redis"
	"TestChat1/model/message"
	"TestChat1/servers/idDispatch"
	"TestChat1/servers/websocket"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type BaseMessage struct {
}

func (this *BaseMessage) CheckSendMessageHasError(msg *message.Message) error {
	return errors.New("禁用此类")
}

//消息放入队列
func (this *BaseMessage) PushMessage(msg *message.Message) error {
	if err := this.CheckSendMessageHasError(msg); err != nil {
		return err
	}
	return this.SelfPushMessage(msg)
}

//自己的东西
func (this *BaseMessage) SelfPushMessage(msg *message.Message) error {
	jsonMessage, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	//发送消息到队列
	defer func() {
		if err != nil {
			fmt.Println("查看一下为什么没关掉")
		}
	}()
	_, err = redis.GoRedisCluster.LPush("message_queue", jsonMessage).Result()
	if err != nil {
		return err
	}

	return nil
}

func (this *BaseMessage) CommonHandle(ml *message.MessageList, resq *sql.Row, msg *message.Message, msgcontent string, isSelf bool, tx *sql.Tx) error {
	resq.Scan(&ml.Id, &ml.MessageNum, &ml.Uid, &ml.FromId)
	if ml.Id == 0 {
		if err := this.InsertData(msg, msgcontent, isSelf, tx); err != nil {
			return err
		}
	} else {
		if err := this.UpdateData(msg, msgcontent, ml.Id, isSelf, tx); err != nil {
			return err
		}
	}
	return nil
}

//查询 我发出去的消息 也就是接收人的消息面板是不是要有红点
func (this *BaseMessage) ISendMessage(ml *message.MessageList, msg *message.Message, msgcontent string, tx *sql.Tx) error {
	resq := tx.QueryRow("SELECT `id`,`message_num`,`uid`,`from_id` FROM `message_list` WHERE `from_id`=? AND `message_type`=? limit 1", msg.SendUid, msg.MessageType)
	return this.CommonHandle(ml, resq, msg, msgcontent, true, tx)
}

//查询 我收的消息里面是不是要有红点
func (this *BaseMessage) IReceiverSendMessage(ml *message.MessageList, msg *message.Message, msgcontent string, tx *sql.Tx) error {
	resq := tx.QueryRow("SELECT `id`,`message_num`,`uid`,`from_id` FROM `message_list` WHERE `uid`=? AND `message_type`=? limit 1", msg.SendUid, msg.MessageType)
	msg.SendUid, msg.ReceiveUid = msg.ReceiveUid, msg.SendUid
	return this.CommonHandle(ml, resq, msg, msgcontent, false, tx)
}

//新增
func (this *BaseMessage) InsertData(msg *message.Message, msgcontent string, isSelf bool, tx *sql.Tx) error {
	//自己发送的消息，不需要红点
	message_num := 0
	if isSelf {
		message_num = 1
	}
	res, err := tx.Exec("INSERT INTO `message_list`"+
		"(`uid`,`from_id`,`message_content`,`message_type`,`created_time`,`update_time`,`message_num`,`chat_id`)"+
		"VALUES (?,?,?,?,?,?,?,?)", msg.ReceiveUid, msg.SendUid, msgcontent, msg.MessageType, msg.CreatedTime, msg.CreatedTime, message_num, msg.Id)
	if err != nil {
		fmt.Println("clientManager line 41", res, err)
		return err
	}
	return nil
}

//修改
func (this *BaseMessage) UpdateData(msg *message.Message, msgcontent string, id int, isSelf bool, tx *sql.Tx) error {
	sqlsql := "UPDATE `message_list` SET " +
		"`message_content`=?,update_time=?,"
	if isSelf {
		sqlsql += "`message_num`=`message_num`+1,"
	}
	sqlsql += "`is_del`=1 WHERE id=?"
	res, err := tx.Exec(sqlsql, msgcontent, msg.CreatedTime, id)
	if err != nil {
		fmt.Println("clientManager line 54", res, err)
		return err
	}
	return nil
}

//目前简单的就只是私发消息的
func (this *BaseMessage) WebSocketRequest(msg *message.Message) {
	c, err := websocket.ClientMangerInstance.GetClient(msg.ReceiveUid)
	if err != nil {
		fmt.Println("错误 --109", err)
		return
	}
	wmsg := common.GetNewWebSocketRequest("GetMessage")
	wmsg.Body = map[string]interface{}{
		"message_content": msg.MessageContent,
		"from_id":         msg.SendUid,
		"message_type":    msg.MessageType,
	}
	c.SendMsg(wmsg)
}

func (this *BaseMessage) GetTitle(longContent string) string {
	//设置可插入list表的消息
	msgTitle := []rune(longContent)
	lmsgTitle := len(msgTitle)
	msgcontent := ""
	if lmsgTitle >= 50 {
		msgcontent = string(msgTitle[0:47]) + "..."
	} else {
		msgcontent = longContent
	}
	return msgcontent
}

//插入message表
func (this *BaseMessage) InsertMessage(msg *message.Message, tx *sql.Tx) error {
	res, err := tx.Exec(
		"INSERT INTO `message`"+
			"(`message_content`,`send_uid`,`receive_uid`,`created_time`,`message_type`,`chat_id`) "+
			"VALUES (?,?,?,?,?,?)", msg.MessageContent, msg.SendUid, msg.ReceiveUid, msg.CreatedTime, msg.MessageType, msg.ChatId)
	if err != nil {
		fmt.Println("clientManager line 139: ", res, err)
		return err
	}
	insertId, err := res.LastInsertId()
	msg.Id = int(insertId)
	return nil
}

//主要逻辑 暂时没什么用
func (this *BaseMessage) AddMessage(msg *message.Message) error {
	//处理消息入库
	tx, err := mysql.DB.Begin()
	msg.CreatedTime = uint64(time.Now().Unix())
	err = this.InsertMessage(msg, tx)
	if err != nil {
		tx.Rollback()
	}
	//获取这个标题
	msgcontent := this.GetTitle(msg.MessageContent)
	ml := new(message.MessageList)
	ml.ChatId = idDispatch.SnowFlakeWorker.GetId()
	if err = this.ISendMessage(ml, msg, msgcontent, tx); err != nil {
		fmt.Println("ClientManager 116 : ", err)
		tx.Rollback()
		return err
	}
	if err = this.IReceiverSendMessage(ml, msg, msgcontent, tx); err != nil {
		fmt.Println("ClientManager 120 : ", err)
		tx.Rollback()
		return err
	}
	tx.Commit()
	msg.ReceiveUid, msg.SendUid = msg.SendUid, msg.ReceiveUid
	//即时发送消息
	go this.WebSocketRequest(msg)
	return nil
}
