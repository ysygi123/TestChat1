package backtask

import (
	"TestChat1/servers/websocket"
	"time"
)

//监控过期的用户，可能出现拔掉网线没有四次挥手，已经短线的客户还在 循环  每0.5s查下一个用户

func MonitoringMain() {
	ticker := time.NewTicker(time.Millisecond * 500)
	for {
		for _, v := range websocket.ClientMangerInstance.Clients {
			select {
			case <-ticker.C:
				websocket.ClientMangerInstance.RWLock.RLock()
				if uint64(time.Now().Unix())-v.HeartBreath > 120 {
					websocket.ClientMangerInstance.CloseChan <- v.Uid
				}
				websocket.ClientMangerInstance.RWLock.RUnlock()
			}
		}
	}
}
