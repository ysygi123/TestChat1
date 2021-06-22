package websocket

import (
	"errors"
	"sync"
)

type WebsocketFunc func(client *Client, message *map[string]string)

var WebSocketRouteManger *WebSocketRoute

type WebSocketRoute struct {
	RWLock *sync.RWMutex
	Route  map[string]WebsocketFunc
}

func WebSocketRouteMangerInit() {
	WebSocketRouteManger = NewWebSocketRoute()
}

func NewWebSocketRoute() *WebSocketRoute {
	return &WebSocketRoute{
		RWLock: new(sync.RWMutex),
		Route:  make(map[string]WebsocketFunc),
	}
}

//注册路由
func (this *WebSocketRoute) RegisterRoute(cmd string, funcName WebsocketFunc) {
	this.RWLock.Lock()
	this.Route[cmd] = funcName
	this.RWLock.Unlock()
}

//获取函数
func (this *WebSocketRoute) GetHandler(cmd string) (h WebsocketFunc, err error) {
	this.RWLock.RLock()
	h, ok := this.Route[cmd]
	if ok == false {
		err = errors.New("查无此功能")
		return nil, err
	}
	return
}
