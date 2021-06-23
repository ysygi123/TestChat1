package messageChild

import (
	"TestChat1/model/message"
	"TestChat1/servers/websocket"
	"fmt"
)

type AddFriendMessage struct {
}

func (this *AddFriendMessage) AddMessage(message *message.PipelineMessage) error {
	fmt.Println(websocket.ClientMangerInstance.GetClient(2))
	return nil
}
