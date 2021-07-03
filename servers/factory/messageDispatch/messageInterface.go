package messageDispatch

import (
	"TestChat1/model/message"
)

type MessageInterface interface {
	//发送消息
	AddMessage(message *message.Message) error
	//推送消息
	PushMessage(message *message.Message) error
}
