package userService

import (
	"TestChat1/common"
	"TestChat1/db/mysql"
	"TestChat1/db/redis"
	"TestChat1/servers/groupService"
	"TestChat1/vaildate/uservalidate"
	"database/sql"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"hash/crc32"
	"strconv"
)

func CheckHadRequestAndHadFriend(sendUid, receiveUid int) error {
	if err := CheckHadFriend(sendUid, receiveUid); err != nil {
		return err
	}
	row, err := mysql.DB.Query("SELECT id FROM `message` WHERE `send_uid`=? AND `receive_uid`=? AND `message_type`=3 LIMIT 1", sendUid, receiveUid)
	if err = CheckHasId(row); err != nil {
		return errors.New("已发起申请")
	}
	return nil
}

func CheckHadFriend(sendUid, receiveUid int) error {
	row, err := mysql.DB.Query("SELECT id FROM `user_friends` WHERE `uid`=? AND `friend_uid`=? LIMIT 1", sendUid, receiveUid)
	if err != nil {
		return err
	}
	if err = CheckHasId(row); err != nil {
		return errors.New("已存在此好友")
	}
	return nil
}

func CheckHasId(row *sql.Rows) error {
	m, err := mysql.GetOneRow(row)
	if err != nil {
		return err
	}
	idStr, ok := m["id"]
	if ok == true && idStr != "0" {
		return errors.New("")
	}
	return nil
}

//获取用户所在的表名
func GetTableName(username string) string {
	crc32Num := crc32.ChecksumIEEE([]byte(username))
	return "user_login_" + strconv.Itoa(int(crc32Num)%10)
}

//判断是否被注册
func CheckHadRegister(username, tableName string) (bool, error) {
	sqlsql := "select id from " + tableName + " where username=?"
	rows, err := mysql.DB.Query(sqlsql, username)
	if err != nil {
		return false, err
	}
	m, err := mysql.GetOneRow(rows)
	if err != nil {
		return false, err
	}
	idStr, ok := m["id"]
	if ok && idStr != "0" {
		return false, nil
	}
	return true, nil
}

func Register(regVal *uservalidate.RegisterValidate) error {
	tableName := GetTableName(regVal.Username)
	hadRegister, err := CheckHadRegister(regVal.Username, tableName)
	if err != nil {
		return err
	}
	if !hadRegister {
		err := errors.New("账号被注册")
		return err
	}
	sqlsql := "insert into " + tableName + " (username,passwd,uid) values (?,?,?)"
	tx, err := mysql.DB.Begin()
	if err != nil {
		return err
	}
	uid, err := getUid()
	if err != nil {
		return err
	}
	_, err = tx.Exec(sqlsql, regVal.Username, common.GetMD5Data(regVal.Passwd), uid)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func Login(loginStruct *uservalidate.LoginValidate) (map[string]string, error) {
	tableName := GetTableName(loginStruct.Username)
	rows, err := mysql.DB.Query("select uid,passwd,username from "+tableName+" where username=? limit 1", loginStruct.Username)
	if err != nil {
		return nil, err
	}
	userData, err := mysql.GetOneRow(rows)
	if err != nil {
		return nil, err
	}
	_, ok := userData["uid"]
	if !ok || userData["uid"] == "0" {
		err := errors.New("无账号")
		return nil, err
	}
	if common.GetMD5Data(loginStruct.Password) != userData["passwd"] {
		err := errors.New("密码错")
		return nil, err
	}
	//获取session
	session := common.GetSession(loginStruct.Username)

	//判断是否登录

	replay, err := redis.GoRedisCluster.Get("uidlogin:" + userData["uid"]).Result()
	if replay != "" {
		err := errors.New("已经登陆")
		return nil, err
	}
	//设置基础信息
	_, err = redis.GoRedisCluster.HMSet(session, map[string]interface{}{"uid": userData["uid"], "username": userData["username"]}).Result()
	if err != nil {
		return nil, err
	}
	//设置是否登录
	_, err = redis.GoRedisCluster.Set("uidlogin:"+userData["uid"], session, 0).Result()
	if err != nil {
		return nil, err
	}
	uid, err := strconv.Atoi(userData["uid"])
	if err != nil {
		return nil, err
	}
	//设置群已经登录。。群多的话会比较慢
	err = SetHasLogin(uid)
	if err != nil {
		return nil, err
	}
	userData["session"] = session
	return userData, nil
}

//返回uid
func getUid() (int, error) {
	i, err := redis.GoRedisCluster.Incr(viper.GetString("uidDisptach")).Result()
	return int(i), err
}

//退出登录
func LoginOut(uid int) (bool, error) {
	uidlogin := fmt.Sprintf(viper.GetString("uidlogin")+"%d", uid)

	session, err := redis.GoRedisCluster.Get(uidlogin).Result()
	if err != nil {
		return false, err
	}
	uidStr, err := redis.GoRedisCluster.HGet(session, "uid").Result()
	if err != nil {
		return false, err
	}
	if uidStr != strconv.Itoa(uid) {
		return false, errors.New("不是这个用户")
	}

	_, err = redis.GoRedisCluster.Del(session).Result()
	_, err = redis.GoRedisCluster.Del(uidlogin).Result()

	groupIds, err := groupService.GetMyGroupId(uid)

	if err != nil {
		return false, err
	}
	for _, groupId := range groupIds {
		groupName := fmt.Sprintf(viper.GetString("groupUser")+"%d", groupId)
		//万一出一个err就很蛋疼 要撤回 先这样放着
		_, err := redis.GoRedisCluster.SRem(groupName, uid).Result()
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

//设置已经登录
func SetHasLogin(uid int) error {
	groupIds, err := groupService.GetMyGroupId(uid)

	if err != nil {
		return err
	}

	for _, groupId := range groupIds {
		groupName := fmt.Sprintf(viper.GetString("groupUser")+"%d", groupId)
		//万一出一个err就很蛋疼 要撤回 先这样放着
		_, err := redis.GoRedisCluster.SAdd(groupName, uid).Result()
		if err != nil {
			return nil
		}
	}

	return nil
}
