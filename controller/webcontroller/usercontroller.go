package webcontroller

import (
	"TestChat1/common"
	"TestChat1/db/mysql"
	"TestChat1/db/redis"
	"TestChat1/model/message"
	"TestChat1/model/user"
	"TestChat1/servers/userService"
	"TestChat1/servers/websocket"
	"TestChat1/vaildate/uservalidate"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func Register(c *gin.Context) {
	var regVal uservalidate.RegisterValidate
	err := common.AutoValidate(c, &regVal)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	err = userService.Register(&regVal)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	common.ReturnResponse(c, 200, 200, "注册成功", nil)
}

func Login(c *gin.Context) {
	var userParams uservalidate.LoginValidate
	err := common.AutoValidate(c, &userParams)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	userData, err := userService.Login(&userParams)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	common.ReturnResponse(c, 200, 200, "登陆成功", map[string]string{
		"session":  userData["session"],
		"uid":      userData["uid"],
		"username": userData["username"],
	})
	return
}

func LookClient(c *gin.Context) {
	websocket.ClientMangerInstance.RWLock.RLock()
	fmt.Printf("%+v\n", websocket.ClientMangerInstance.Clients)
	for c, v := range websocket.ClientMangerInstance.Clients {
		fmt.Println("啊哈哈哈", c, v)
	}
	websocket.ClientMangerInstance.RWLock.RUnlock()
}

//修改client的字段 后续要改成用redis队列丢进去修改将websocket和web拆开 作废
func AuthClient(c *gin.Context) {
	authParams := uservalidate.Auth{}
	err := common.AutoValidate(c, &authParams)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	rec := redis.RedisPool.Get()
	defer rec.Close()
	reply, err := rec.Do("HGET", authParams.Session, "uid")
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	if reply == nil {
		common.ReturnResponse(c, 200, 400, "未知错误", nil)
		return
	}
	uidStr := string(reply.([]byte))
	uid, _ := strconv.Atoi(uidStr)
	if uid != authParams.Uid {
		common.ReturnResponse(c, 200, 400, "去你妈的吧 uid不对等", nil)
		return
	}
	err = websocket.ClientMangerInstance.SetAuth(uid)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	common.ReturnResponse(c, 200, 200, "成功", nil)
}

//加好友请求 暂时作废
func AddFriendRequest(c *gin.Context) {
	userAddRequest := &uservalidate.AddFriendRequest{}
	err := common.AutoValidate(c, userAddRequest)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	rec := redis.RedisPool.Get()
	defer rec.Close()
	msg := &message.Message{
		MessageType:    uint8(3),
		SendUid:        userAddRequest.SendUid,
		ReceiveUid:     userAddRequest.ReceiveUid,
		CreatedTime:    uint64(time.Now().Unix()),
		MessageContent: userAddRequest.Rname + "向您发出好友请求",
	}
	err = userService.CheckHadFriend(msg.SendUid, msg.ReceiveUid)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
	}
	jsonStr, err := json.Marshal(msg)
	_, err = rec.Do("LPUSH", "message_queue", jsonStr)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	common.ReturnResponse(c, 200, 200, "success", nil)
}

//同意加好友
func AddFriendCommit(c *gin.Context) {
	addFriendCommit := &uservalidate.AddFriendCommit{}
	err := common.AutoValidate(c, addFriendCommit)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	msg := &message.Message{}

	row := mysql.DB.QueryRow("SELECT `send_uid`,`receive_uid` FROM `message` WHERE `id`=?", addFriendCommit.MessageId)
	t := uint64(time.Now().Unix())
	err = row.Scan(&msg.SendUid, &msg.ReceiveUid)
	err = userService.CheckHadFriend(msg.SendUid, msg.ReceiveUid)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}

	tx, err := mysql.DB.Begin()
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	_, err = tx.Exec("INSERT INTO `user_friends` (`uid`,`friend_uid`,`created_time`,`update_time`) VALUES "+
		"(?,?,?,?),(?,?,?,?)",
		msg.SendUid, msg.ReceiveUid, t, t, msg.ReceiveUid, msg.SendUid, t, t)
	if err != nil {
		common.ReturnResponse(c, 200, 400, "已经存在此好友", nil)
		return
	}
	_, err = tx.Exec("insert into message_list  (uid,from_id,message_content,message_type,created_time,update_time,message_num) values "+
		"(?,?,'你们已经成为好友',1,?,?,1),"+
		"(?,?,'你们已经成为好友',1,?,?,1)",
		msg.SendUid, msg.ReceiveUid, t, t, msg.ReceiveUid, msg.SendUid, t, t)
	tx.Commit()
	if err != nil {
		common.ReturnResponse(c, 200, 400, "已经存在此好友", nil)
		return
	}
	common.ReturnResponse(c, 200, 200, "success", nil)
}

//获取好友列表
func GetFriendsList(c *gin.Context) {
	uidStr := c.Param("uid")
	uid, _ := strconv.Atoi(uidStr)
	if uid == 0 {
		common.ReturnResponse(c, 200, 400, "id错误", nil)
		return
	}
	rows, err := mysql.DB.Query("SELECT `u2`.`username`,`u2`.`rname`,`u2`.`uid`,`u2`.`mobile` "+
		"FROM `user_friends` as `u1` INNER JOIN `user` as `u2` ON `u1`.`uid`=`u2`.`uid`"+
		"WHERE `u1`.`uid`=? AND is_del=1", uid)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	userList := make([]user.User, 0)
	var tmpUser user.User
	for rows.Next() {
		err = rows.Scan(&tmpUser.Username, &tmpUser.Rname, &tmpUser.Uid, &tmpUser.Mobile)
		if err != nil {
			common.ReturnResponse(c, 200, 400, err.Error(), nil)
			return
		}
		userList = append(userList, tmpUser)
	}
	common.ReturnResponse(c, 200, 200, "成功", map[string]interface{}{"list": userList})
}
