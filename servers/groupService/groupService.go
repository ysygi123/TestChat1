package groupService

import (
	"TestChat1/db/mysql"
	"errors"
)

func CheckHadGroup(groupId, uid int) error {
	row, err := mysql.DB.Query("SELECT `id` FROM `group_user` WHERE `uid`=? AND `group_id`=?", uid, groupId)
	if err != nil {
		return err
	}
	m, err := mysql.GetOneRow(row)
	if err != nil {
		return err
	}
	id, ok := m["id"]
	if id != "0" && ok {
		return errors.New("已经加入该群")
	}
	return nil
}
