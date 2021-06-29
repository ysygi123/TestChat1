package messageChild

import (
	"TestChat1/common"
	"TestChat1/db/mysql"
	"TestChat1/db/redis"
	"TestChat1/model/message"
	"TestChat1/servers/websocket"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type GroupMessage struct {
	BaseMessage
}

func (this *GroupMessage) CheckSendMessageHasError(msg *message.Message) error {
	if msg.GroupId == 0 || msg.ReceiveUid != 0 {
		return errors.New("群消息格式错误 group_id : " + strconv.Itoa(msg.GroupId) + "  receiveuid : " + strconv.Itoa(msg.ReceiveUid))
	}
	return nil
}

func (this *GroupMessage) PushMessage(msg *message.Message) error {
	if err := this.CheckSendMessageHasError(msg); err != nil {
		return err
	}
	return this.SelfPushMessage(msg)
}

func (this *GroupMessage) AddMessage(msg *message.Message) error {
	allUids, err := this.getThisGroupUserIds(msg.GroupId)
	if err != nil {
		return err
	}
	allOnLineUids := this.getIsLoginUids(allUids)
	if err := this.SetInDataBase(allUids, msg); err != nil {
		return err
	}
	go this.WebSocketRequest(msg, allOnLineUids)

	return nil
}

//向这些人发送websocket
func (this *GroupMessage) WebSocketRequest(msg *message.Message, uids []int) {
	clients, err := websocket.ClientMangerInstance.GetManyClient(uids)
	if err != nil {
		fmt.Println("错误 群发消息48")
		return
	}
	wmsg := common.GetNewWebSocketRequest("GetMessage")
	wmsg.Body = map[string]interface{}{
		"message_content": msg.MessageContent,
		"from_id":         msg.SendUid,
		"message_type":    msg.MessageType,
		"group_id":        msg.GroupId,
	}
	//后面要改 通过channel发给每个*client.WebsocketConn 各自开启一个协程阻塞监听
	for _, c := range clients {
		c.SendMsg(wmsg)
	}
}

//数据库操作
func (this *GroupMessage) SetInDataBase(allUids []int, msg *message.Message) error {
	mysql.DB.Exec("begin")
	if err := this.InsertMessage(msg); err != nil {
		return err
	}
	title := this.GetTitle(msg.MessageContent)
	fmt.Println("查看群聊的title是什么 : ", title)
	sqlsql := "update message_list set " +
		"message_content=? and message_num=message_num+1 and update_time=? and is_del=1 where from_id=? and message_type=2 and uid in (" +
		common.IntJoin(allUids, len(allUids)) + ")"
	_, err := mysql.DB.Query(sqlsql, title, uint64(time.Now().Unix()), msg.GroupId)
	if err != nil {
		return err
	}
	//默认加群就给一条消息，省去新增步骤
	mysql.DB.Exec("commit")
	return nil
}

//获取这个群里的所有uid
func (this *GroupMessage) getThisGroupUserIds(groupId int) ([]int, error) {
	rows, err := mysql.DB.Query("select uid from group_users where group_id=?", groupId)
	if err != nil {
		return nil, err
	}
	var uidsSlice []int
	var tmpUid int

	for rows.Next() {
		rows.Scan(&tmpUid)
		uidsSlice = append(uidsSlice, tmpUid)
	}

	return append(uidsSlice, []int{}...), nil
}

//获取群里已经登录的uid
func (this *GroupMessage) getIsLoginUids(uids []int) []int {
	strKeys := make([]interface{}, 0)
	for _, uid := range uids {
		strKeys = append(strKeys, "uidlogin:"+strconv.Itoa(uid))
	}
	rec := redis.RedisPool.Get()
	defer rec.Close()
	replay, err := rec.Do("MGET", strKeys...)
	if err != nil {
		return nil
	}
	returnInt := make([]int, 0)
	re := replay.([]interface{})
	for i, v := range re {
		if v != nil {
			returnInt = append(returnInt, uids[i])
		}
	}
	return returnInt
}
