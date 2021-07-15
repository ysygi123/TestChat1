package backtask

import (
	"TestChat1/db/redis"
	"TestChat1/model/message"
	"TestChat1/servers/factory/messageDispatch"
	"TestChat1/servers/websocket"
	"encoding/json"
	"fmt"
	"sync"
)

//后台任务消费消息
func TaskConsumeMessage() {
	wg := new(sync.WaitGroup)
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			for {
				reply, err := redis.GoRedisCluster.BRPop(0, "message_queue").Result()
				if err != nil {
					fmt.Println("clientManager line 84: ", err)
					continue
				}
				msg := new(message.Message)
				err = json.Unmarshal([]byte(reply[1]), msg)
				if err != nil {
					fmt.Println("clientManager line 90: ", err)
					continue
				}

				mfc, err := messageDispatch.CreateMessage(map[string]interface{}{
					"messageType": msg.MessageType,
				})
				if err != nil {
					fmt.Println("clientManager line 40", err)
					continue
				}
				err = mfc.AddMessage(msg)
				if err != nil {
					fmt.Println("clientManager line 44", err)
					continue
				}
			}
		}(wg)
	}
	wg.Done()
}

//接收删除客户fd的请求
func CleanClient() {
	for {
		c := <-websocket.ClientMangerInstance.CloseChan
		websocket.ClientMangerInstance.DelClient(c)
		fmt.Println("我删除了", c)
	}
}
