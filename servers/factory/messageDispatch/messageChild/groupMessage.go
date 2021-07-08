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
	_, err = tx.Exec(sqlsql, uint64(time.Now().Unix()), msg.GroupId)
	if err != nil {
		return err
	}
	//默认加群就给一条消息，省去新增步骤
	tx.Commit()
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
	returnInt := make([]int, 0)
	for _, uid := range uids {
		replay, err := redis.GoRedisCluster.Get("uidlogin:" + strconv.Itoa(uid)).Result()
		if err != nil {
			continue
		}
		if replay != "" {
			returnInt = append(returnInt, uid)
		}
	}
	return returnInt
}
