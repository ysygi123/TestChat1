package messageChild

import (
	"TestChat1/model/message"
	"errors"
)

type AddFriendMessage struct {
	BaseMessage
}

func (this *AddFriendMessage) CheckSendMessageHasError(msg *message.Message) error {
	if msg.GroupId != 0 {
		return errors.New("禁止发送群消息")
	}
	return nil
}

func (this *AddFriendMessage) PushMessage(msg *message.Message) error {
	if err := this.CheckSendMessageHasError(msg); err != nil {
		return err
	}
	return this.SelfPushMessage(msg)
}
