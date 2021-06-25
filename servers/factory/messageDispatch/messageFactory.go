package messageDispatch

import (
	"TestChat1/servers/factory/messageDispatch/messageChild"
	"errors"
	"fmt"
)

type MessageOpsFactory func(conf map[string]interface{}) (MessageInterface, error)

var MessageFactory = make(map[uint8]MessageOpsFactory)

func init() {
	Register(uint8(0), NewAddBaseMessage)
	Register(uint8(1), NewUserMessage)
	Register(uint8(3), NewAddFriendMessage)
}

//注册
func Register(messageType uint8, factory MessageOpsFactory) {
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
		err := errors.New("没有这个类1")
		return nil, err
	}
	opsFactory, ok := MessageFactory[opsType.(uint8)]
	if !ok {
		err := errors.New("没有这个类2")
		return nil, err
	}
	return opsFactory(conf)
}

//new一个的个人消息
func NewUserMessage(conf map[string]interface{}) (MessageInterface, error) {
	fmt.Println("create usermessagecreate")
	return &messageChild.UserMessage{}, nil
}

func NewAddFriendMessage(conf map[string]interface{}) (MessageInterface, error) {
	fmt.Println("CREATE AddFriendmessage")
	return &messageChild.AddFriendMessage{}, nil
}

func NewAddBaseMessage(conf map[string]interface{}) (MessageInterface, error) {
	fmt.Println("create baseMessage")
	return &messageChild.BaseMessage{}, nil
}
