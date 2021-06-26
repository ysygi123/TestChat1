package userService

import (
	"TestChat1/db/mysql"
	"database/sql"
	"errors"
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
