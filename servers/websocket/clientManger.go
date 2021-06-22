package websocket

import (
	"errors"
	"sync"
)

type ClientManger struct {
	Clients map[int]*Client
	RWLock  *sync.RWMutex
}

func NewClientManger() *ClientManger {
	return &ClientManger{
		Clients: map[int]*Client{},
		RWLock:  new(sync.RWMutex),
	}
}

func (this *ClientManger) AddClient(uid int, client *Client) {
	this.RWLock.Lock()
	this.Clients[uid] = client
	this.RWLock.Unlock()
}

func (this *ClientManger) GetClient(uid int) (c *Client, e error) {
	this.RWLock.RLock()
	c, ok := this.Clients[uid]
	this.RWLock.RUnlock()
	if ok == false {
		e = errors.New("没有这个用户")
		return nil, e
	}
	return
}

func (this *ClientManger) DelClient(uid int) {
	this.RWLock.Lock()
	delete(this.Clients, uid)
	this.RWLock.Unlock()
}
