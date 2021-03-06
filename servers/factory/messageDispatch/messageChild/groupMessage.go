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
	if msg.ChatId == 0 || msg.ReceiveUid != 0 {
		return errors.New("群消息格式错误 receiveuid : " + strconv.Itoa(msg.ReceiveUid))
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
	allUids, err := this.getThisGroupUserIds(msg.ChatId)
	if err != nil {
		return err
	}
	allOnLineUids := this.getIsLoginUids(msg.ChatId)
	go this.WebSocketRequest(msg, allOnLineUids)
	if err := this.setInDataBase(allUids, msg); err != nil {
		return err
	}

	return nil
}

//向这些人发送websocket
func (this *GroupMessage) WebSocketRequest(msg *message.Message, uids []int) {
	clients, err := websocket.ClientMangerInstance.GetManyClient(uids, msg.SendUid)
	if err != nil {
		fmt.Println("错误 群发消息48")
		return
	}
	wmsg := common.GetNewWebSocketRequest("GetMessage")
	wmsg.Body = map[string]interface{}{
		"message_content": msg.MessageContent,
		"from_id":         msg.SendUid,
		"message_type":    msg.MessageType,
		"chat_id":         msg.ChatId,
	}
	//后面要改 通过channel发给每个*client.WebsocketConn 各自开启一个协程阻塞监听
	for _, c := range clients {
		c.SendMsg(wmsg)
	}
}

//数据库操作
func (this *GroupMessage) setInDataBase(allUids []int, msg *message.Message) error {
	tx, err := mysql.DB.Begin()
	if err != nil {
		return err
	}
	if err := this.InsertMessage(msg, tx); err != nil {
		return err
	}
	title := this.GetTitle(msg.MessageContent)
	sqlsql := "update message_list set message_content='" + title + "',message_num=message_num+1,update_time=?,is_del=1 where from_id=? and message_type=2 and uid in (" +
		common.IntJoin(allUids, len(allUids)) + ")"
	_, err = tx.Exec(sqlsql, uint64(time.Now().Unix()), msg.ChatId)
	if err != nil {
		return err
	}
	//默认加群就给一条消息，省去新增步骤
	tx.Commit()
	return nil
}

//获取这个群里的所有uid
func (this *GroupMessage) getThisGroupUserIds(groupId uint64) ([]int, error) {
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
func (this *GroupMessage) getIsLoginUids(groupId uint64) []int {
	allOnlineUid := make([]int, 0)
	uidStrs, err := redis.GoRedisCluster.SMembers(fmt.Sprintf("group_user:%d", groupId)).Result()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, uidStr := range uidStrs {
		uid, _ := strconv.Atoi(uidStr)
		allOnlineUid = append(allOnlineUid, uid)
	}
	return allOnlineUid
}
