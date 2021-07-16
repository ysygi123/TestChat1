package messageChild

import (
	"TestChat1/db/mysql"
	"TestChat1/model/message"
	"TestChat1/servers/idDispatch"
	"TestChat1/servers/userService"
	"database/sql"
	"errors"
	"time"
)

type AddFriendMessage struct {
	BaseMessage
}

func (this *AddFriendMessage) CheckSendMessageHasError(msg *message.Message) error {
	if msg.GroupId != 0 {
		return errors.New("禁止发送群消息")
	}
	if msg.ReceiveUid <= 0 || msg.SendUid <= 0 {
		return errors.New("好友id错误")
	}
	msg.MessageContent = "你有一个好友请求"
	return nil
}

func (this *AddFriendMessage) PushMessage(msg *message.Message) error {
	if err := this.CheckSendMessageHasError(msg); err != nil {
		return err
	}
	err := userService.CheckHadRequestAndHadFriend(msg.SendUid, msg.ReceiveUid)
	if err != nil {
		return err
	}
	return this.SelfPushMessage(msg)
}

//主要逻辑 发送好友请求
func (this *AddFriendMessage) AddMessage(msg *message.Message) error {
	//处理消息入库
	tx, err := mysql.DB.Begin()
	if err != nil {
		return err
	}
	row := tx.QueryRow("select id,chat_id from message_list where uid=? and message_type=3 limit 1")
	var id int
	var chatId uint64
	err = row.Scan(&id, &chatId)
	//err有错 就认为是这条位空
	if err != nil {
		msg.ChatId = idDispatch.SnowFlakeWorker.GetId()
		_, err := tx.Exec("insert into message_list (uid,message_content,message_type,created_time,update_time,message_num, chat_id) values (?,?,3,?,?,1,?)",
			msg.ReceiveUid, msg.MessageContent, msg.CreatedTime, msg.CreatedTime, msg.ChatId)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		msg.ChatId = chatId
		_, err := tx.Exec("update message_list set update_time,message_num=message_num+1,is_del=1 where uid=? and message_type=3 limit 1", msg.ReceiveUid)
		if err != nil {
			tx.Exec("rollback")
			return err
		}
	}
	msg.CreatedTime = uint64(time.Now().Unix())
	err = this.InsertMessage(msg, tx)
	if err != nil {
		tx.Rollback()
		return nil
	}

	//即时发送消息
	go this.WebSocketRequest(msg)
	return nil
}

func (this *AddFriendMessage) InsertMessage(msg *message.Message, tx *sql.Tx) error {
	_, err := tx.Exec("insert into message (send_uid, receive_uid, created_time,message_type,chat_id) values (?,?,?,3,?)",
		msg.SendUid, msg.ReceiveUid, msg.CreatedTime, msg.ChatId)
	if err != nil {
		return err
	}
	return nil
}
