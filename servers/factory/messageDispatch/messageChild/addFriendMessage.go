package messageChild

import (
	"TestChat1/db/mysql"
	"TestChat1/model/message"
	"TestChat1/servers/userService"
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
	_, err := mysql.DB.Exec("BEGIN")
	msg.CreatedTime = uint64(time.Now().Unix())
	err = this.InsertMessage(msg)
	if err != nil {
		mysql.DB.Exec("ROLLBACK")
		return nil
	}

	rows, err := mysql.DB.Query("select id from message_list where uid=? and message_type=3 limit 1")
	if err != nil {
		return err
	}
	m, err := mysql.GetOneRow(rows)
	if err != nil {
		mysql.DB.Exec("rollback")
		return err
	}

	idStr, ok := m["id"]

	if idStr == "0" || ok == false {
		_, err := mysql.DB.Exec("insert into message_list (uid,message_content,message_type,created_time,update_time,message_num) values (?,?,3,?,?,1)",
			msg.ReceiveUid, msg.MessageContent, msg.CreatedTime, msg.CreatedTime)
		if err != nil {
			mysql.DB.Exec("rollback")
			return err
		}
	} else {
		_, err := mysql.DB.Exec("update message_list set update_time,message_num=message_num+1,is_del=1 where uid=? and message_type=3 limit 1", msg.ReceiveUid)
		if err != nil {
			mysql.DB.Exec("rollback")
			return err
		}
	}
	//即时发送消息
	go this.WebSocketRequest(msg)
	return nil
}

func (this *AddFriendMessage) InsertMessage(msg *message.Message) error {
	_, err := mysql.DB.Exec("insert into SetInDataBase (send_uid, receive_uid, created_time,message_type) values (?,?,?,3)",
		msg.SendUid, msg.ReceiveUid, msg.CreatedTime)
	if err != nil {
		return err
	}
	return nil
}
