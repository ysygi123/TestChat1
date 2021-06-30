package messageChild

import (
	"TestChat1/db/mysql"
	"TestChat1/model/message"
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

//主要逻辑
func (this *AddFriendMessage) AddMessage(msg *message.Message) error {
	//处理消息入库
	tx, err := mysql.DB.Begin()
	if err != nil {
		return err
	}
	msg.CreatedTime = uint64(time.Now().Unix())
	err = this.InsertMessage(msg, tx)
	if err != nil {
		tx.Rollback()
		return nil
	}

	rows, err := tx.Query("select id from message_list where uid=? and message_type=3 limit 1")
	if err != nil {
		return err
	}
	m, err := mysql.GetOneRow(rows)
	if err != nil {
		tx.Rollback()
		return err
	}

	idStr, ok := m["id"]

	if idStr == "0" || ok == false {
		_, err := tx.Exec("insert into message_list (uid,message_content,message_type,created_time,update_time,message_num) values (?,?,3,?,?,1)",
			msg.ReceiveUid, msg.MessageContent, msg.CreatedTime, msg.CreatedTime)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		_, err := tx.Exec("update message_list set update_time,message_num=message_num+1,is_del=1 where uid=? and message_type=3 limit 1", msg.ReceiveUid)
		if err != nil {
			tx.Exec("rollback")
			return err
		}
	}
	//即时发送消息
	go this.WebSocketRequest(msg)
	return nil
}

func (this *AddFriendMessage) InsertMessage(msg *message.Message, tx *sql.Tx) error {
	_, err := tx.Exec("insert into SetInDataBase (send_uid, receive_uid, created_time,message_type) values (?,?,?,3)",
		msg.SendUid, msg.ReceiveUid, msg.CreatedTime)
	if err != nil {
		return err
	}
	return nil
}
