package websocket

import (
	"TestChat1/common"
	"TestChat1/db/mysql"
	"TestChat1/db/redis"
	"TestChat1/model/message"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"strconv"
	"sync"
	"time"
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

func (this *ClientManger) DelClient(uid int) {
	this.RWLock.Lock()
	delete(this.Clients, uid)
	this.RWLock.Unlock()
}

func (this *ClientManger) LoopToKillChild() {
	for {
		x := <-this.CloseChan
		this.DelClient(x)
	}
}

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

//异步消费这些垃圾
func (this *ClientManger) ConsumeMessage() {
	wg := new(sync.WaitGroup)
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			rec := redis.RedisPool.Get()
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
				go addMessage(msg)
				c, err := this.GetClient(msg.ReceiveUid)
				if err != nil {
					fmt.Println("clientManager line 100: ", err)
					continue
				}
				go func() {
					m := common.GetNewWebSocketRequest("GetMsg")
					m.Body["message"] = msg.MessageContent
					err = c.WebSocketConn.WriteMessage(websocket.TextMessage, common.GetJsonByteData(m))
					if err != nil {
						fmt.Println("clientManager line 110: ", err)
					}
				}()
			}
		}(wg)
	}
	wg.Done()
}

func addMessage(msg *message.Message) {
	res, err := mysql.DB.Exec("BEGIN")
	msg.CreatedTime = uint64(time.Now().Unix())
	res, err = mysql.DB.Exec(
		"INSERT INTO `message`"+
			"(`message_content`,`send_uid`,`receive_uid`,`created_time`) "+
			"VALUES (?,?,?,?)", msg.MessageContent, msg.SendUid, msg.ReceiveUid, msg.CreatedTime)
	if err != nil {
		fmt.Println("clientManager line 126: ", res, err)
		mysql.DB.Exec("ROLLBACK")
	}

	resq, err := mysql.DB.Query("SELECT `id`,`message_num` FROM `message_list` WHERE `from_id`=? AND `message_type`=1 AND `is_del`=1 limit 1", msg.SendUid)
	if err != nil {
		fmt.Println("clientManager line 131", resq, err)
	}
	resqMap, err := mysql.GetOneRow(resq)
	//设置可插入list表的消息
	msgTitle := []rune(msg.MessageContent)
	lmsgTitle := len(msgTitle)
	msgcontent := ""
	if lmsgTitle >= 50 {
		msgcontent = string(msgTitle[0:47]) + "..."
	} else {
		msgcontent = msg.MessageContent
	}

	if idStr, ok := resqMap["id"]; ok == false { //新增
		res, err := mysql.DB.Exec("INSERT INTO `message_list`"+
			"(`uid`,`from_id`,`message_content`,`message_type`,`created_time`,`update_time`,`message_num`)"+
			"VALUES (?,?,?,?,?,?,?)", msg.ReceiveUid, msg.SendUid, msgcontent, 1, msg.CreatedTime, msg.CreatedTime, 1)
		if err != nil {
			fmt.Println("clientManager line 150", res, err)
			mysql.DB.Exec("ROLLBACK")
		}
	} else { //修改
		id, _ := strconv.Atoi(idStr)
		res, err := mysql.DB.Exec("UPDATE `message_list` SET "+
			"`message_content`=?,update_time=?,`message_num`=`message_num`+1 "+
			"WHERE id=?",
			msgcontent, msg.CreatedTime, id)
		if err != nil {
			fmt.Println("clientManager line 159", res, err)
			mysql.DB.Exec("ROLLBACK")
		}
	}
	res, err = mysql.DB.Exec("COMMIT")
}
