package messageChild

import (
	"TestChat1/model/message"
	"errors"
)

type UserMessage struct {
	BaseMessage
}

func (this *UserMessage) CheckSendMessageHasError(msg *message.Message) error {
	if msg.GroupId != 0 {
		return errors.New("禁止发送群消息")
	}
	return nil
}

func (this *UserMessage) PushMessage(msg *message.Message) error {
	if err := this.CheckSendMessageHasError(msg); err != nil {
		return err
	}
	return this.SelfPushMessage(msg)
}
