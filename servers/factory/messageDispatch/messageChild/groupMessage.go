package messageChild

import (
	"TestChat1/model/message"
	"errors"
)

type GroupMessage struct {
	BaseMessage
}

func (this *GroupMessage) CheckSendMessageHasError(msg *message.Message) error {
	if msg.GroupId == 0 || msg.ReceiveUid != 0 {
		return errors.New("群消息格式错误")
	}
	return nil
}

func (this *GroupMessage) PushMessage(msg *message.Message) error {
	if err := this.CheckSendMessageHasError(msg); err != nil {
		return err
	}
	return this.SelfPushMessage(msg)
}

func (this *GroupMessage) SelfPushMessage(msg *message.Message) error {

	return nil
}

func (this *GroupMessage) AddMessage(msg *message.Message) error {
	return nil
}
