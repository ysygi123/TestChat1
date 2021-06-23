package messageDispatch

import (
	"TestChat1/servers/factory/messageDispatch/messageChild"
	"errors"
	"fmt"
)

type MessageOpsFactory func(conf map[string]interface{}) (MessageInterface, error)

var MessageFactory = make(map[string]MessageOpsFactory)

func init() {
	Register("1", NewUserMessage)
}

//注册
func Register(messageType string, factory MessageOpsFactory) {
	if factory == nil {
		fmt.Println("没有传啊")
		return
	}
	_, ok := MessageFactory[messageType]
	if ok {
		fmt.Println("已经存在")
		return
	} else {
		MessageFactory[messageType] = factory
	}
}

func CreateMessage(conf map[string]interface{}) (MessageInterface, error) {
	opsType, ok := conf["messageType"]
	if !ok {
		err := errors.New("没有这个类")
		return nil, err
	}
	opsFactory, ok := MessageFactory[opsType.(string)]
	if !ok {
		err := errors.New("没有这个类")
		return nil, err
	}
	return opsFactory(conf)
}

//new一个的个人消息
func NewUserMessage(conf map[string]interface{}) (MessageInterface, error) {
	fmt.Println("usermessagecreate")
	return &messageChild.UserMessage{}, nil
}
