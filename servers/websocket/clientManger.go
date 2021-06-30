package websocket

import (
	"errors"
	"fmt"
	"sync"
)

type ClientManger struct {
	Clients   map[int]*Client
	RWLock    *sync.RWMutex
	CloseChan chan int //关闭通道收到则删除某个客户
}

var (
	ClientMangerInstance *ClientManger
)

func ClientMangerInstanceInit() {
	ClientMangerInstance = NewClientManger()
}

func NewClientManger() *ClientManger {
	return &ClientManger{
		Clients:   map[int]*Client{},
		RWLock:    new(sync.RWMutex),
		CloseChan: make(chan int, 1000),
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

//批量获取Client的链接
func (this *ClientManger) GetManyClient(uids []int) (c []*Client, e error) {
	this.RWLock.RLock()
	for k, v := range uids {
		if (k+1)%100 == 0 {
			this.RWLock.RUnlock()
			this.RWLock.RLock()
		}
		cc, ok := this.Clients[v]
		if !ok {
			continue
		}
		c = append(c, cc)
	}
	this.RWLock.RUnlock()
	return
}

func (this *ClientManger) DelClient(uid int) {
	this.RWLock.Lock()
	defer this.RWLock.Unlock()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r, "啊我去")
		}
	}()
	delete(this.Clients, uid)
}

func (this *ClientManger) LoopToKillChild() {
	for {
		x := <-this.CloseChan
		this.DelClient(x)
	}
}

//设置一个标示，标示这个人已经经过认证
func (this *ClientManger) SetAuth(uid int) (err error) {
	c, err := this.GetClient(uid)
	if err != nil {
		return
	}
	this.RWLock.Lock()

	c.IsAuth = true
	this.RWLock.Unlock()
	return
}
