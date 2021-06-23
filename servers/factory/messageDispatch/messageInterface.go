package messageDispatch

import "TestChat1/model/message"

type MessageInterface interface {
	AddMessage(message *message.PipelineMessage) error
}
