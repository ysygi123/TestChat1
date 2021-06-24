package messageDispatch

import (
	"TestChat1/model/message"
)

type MessageInterface interface {
	AddMessage(message *message.Message) error
	PushMessage(message *message.Message) error
}
