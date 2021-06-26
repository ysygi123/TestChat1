package messageChild

import (
	"TestChat1/model/message"
	"TestChat1/servers/UserService"
	"errors"
)

type AddFriendMessage struct {
	BaseMessage
}

func (this *AddFriendMessage) CheckSendMessageHasError(msg *message.Message) error {
	if msg.GroupId != 0 {
		return errors.New("禁止发送群消息")
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
