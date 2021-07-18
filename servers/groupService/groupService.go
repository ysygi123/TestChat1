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

//查找这个人有多少个群
func GetMyGroupId(uid int) ([]int, error) {
	rows, err := mysql.DB.Query("select group_id from group_users where uid=? and is_del=1", uid)
	if err != nil {
		return nil, err
	}
	tmpId := 0
	groupIds := make([]int, 0)
	for rows.Next() {
		rows.Scan(&tmpId)
		groupIds = append(groupIds, tmpId)
	}
	return groupIds, nil
}
