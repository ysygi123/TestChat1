package messageChild

import (
	"TestChat1/db/mysql"
	"TestChat1/model/message"
	"errors"
	"time"
)

type UserMessage struct {
	BaseMessage
}

func (this *UserMessage) CheckSendMessageHasError(msg *message.Message) error {
	if msg.ChatId == 0 {
		return errors.New("chatid不为0")
	}
	return nil
}

func (this *UserMessage) PushMessage(msg *message.Message) error {
	if err := this.CheckSendMessageHasError(msg); err != nil {
		return err
	}
	return this.SelfPushMessage(msg)
}

//消息入库
func (this *UserMessage) AddMessage(msg *message.Message) error {
	//处理消息入库
	tx, err := mysql.DB.Begin()
	msg.CreatedTime = uint64(time.Now().Unix())
	err = this.InsertMessage(msg, tx)
	if err != nil {
		tx.Rollback()
	}
	//获取这个标题
	msgcontent := this.GetTitle(msg.MessageContent)
	_, err = tx.Exec("update message_list set message_content=?,is_del=1 where chat_id=? and message_type=1", msgcontent, msg.ChatId)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("update message_list set message_num=message_num+1 where chat_id=? and message_type=1 and uid=?", msg.ChatId, msg.ReceiveUid)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	//即时发送消息
	go this.WebSocketRequest(msg)
	return nil
}
