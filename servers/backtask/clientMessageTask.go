package backtask

import (
	"TestChat1/db/redis"
	"TestChat1/model/message"
	"TestChat1/servers/factory/messageDispatch"
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
			rec := redis.RedisPool.Get()
			defer rec.Close()
			for {
				reply, err := rec.Do("BRPOP", "message_queue", 0)
				if err != nil {
					fmt.Println("clientManager line 84: ", err)
					continue
				}
				msg := new(message.Message)
				err = json.Unmarshal(reply.([]interface{})[1].([]byte), msg)
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
