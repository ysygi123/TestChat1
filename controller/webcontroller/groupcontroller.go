package webcontroller

import (
	"TestChat1/common"
	"TestChat1/db/mysql"
	"TestChat1/vaildate/groupvalidate"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

//创建群
func CreateGroup(c *gin.Context) {
	groupCreateValidate := &groupvalidate.GroupCreateValidate{}
	if err := common.AutoValidate(c, groupCreateValidate); err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	t := uint64(time.Now().Unix())
	tx, err := mysql.DB.Begin()
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	result, err := tx.Exec("INSERT INTO `group` (`group_name`,`created_uid`,`created_time`,`update_time`) VALUES (?,?,?,?)",
		groupCreateValidate.GroupName, groupCreateValidate.CreatedUid, t, t)
	if err != nil {
		tx.Rollback()
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	groupId, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	_, err = tx.Exec("INSERT INTO `group_users` (`uid`,`group_id`, `created_time`,`update_time`) VALUES (?,?,?,?)",
		groupCreateValidate.CreatedUid, groupId, t, t)
	if err != nil {
		tx.Rollback()
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	tx.Commit()
	common.ReturnResponse(c, 200, 200, "success", nil)
}

func GroupList(c *gin.Context) {
	uidStr := c.Param("uid")
	uid, _ := strconv.Atoi(uidStr)
	if uid == 0 {
		common.ReturnResponse(c, 200, 400, "id错误", nil)
		return
	}
	rows, err := mysql.DB.Query("SELECT `g.id`,`g.group_name` FROM `group_user as gu` LEFT JOIN `group as g` " +
		"ON `gu.group_id`=`g.id` WHERE `g`.`is_del`=1 AND `gu`.`is_del`=1")
	if err != nil {
		common.ReturnResponse(c, 200, 400, "id错误", nil)
		return
	}
	list, err := mysql.GetManyRows(rows)
	if err != nil {
		common.ReturnResponse(c, 200, 400, "id错误", nil)
		return
	}
	common.ReturnResponse(c, 200, 400, "success", list)
}

//加入群
func AddToGroupCommit(c *gin.Context) {
	addToGroupCommit := &groupvalidate.AddToGroupCommitValidate{}
	if err := common.AutoValidate(c, addToGroupCommit); err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	tx, err := mysql.DB.Begin()
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	t := uint64(time.Now().Unix())
	_, err = tx.Exec("insert into group_users (uid,group_id,created_time,update_time) values (?,?,?,?)",
		addToGroupCommit.Uid, addToGroupCommit.GroupId, t, t)
	if err != nil {
		tx.Rollback()
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}

	_, err = tx.Exec("insert into message (message_content,created_time,chat_id,message_type) "+
		"values ('欢迎加入群',?,?,?)", t, addToGroupCommit.GroupId, 2)
	if err != nil {
		tx.Rollback()
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}

	_, err = tx.Exec("insert into message_list(message_content,uid,from_id,message_type,created_time,update_time,message_num,chat_id) values ("+
		"'欢迎加入群',?,?,2,?,?,1,?)",
		addToGroupCommit.Uid, addToGroupCommit.GroupId, t, t, addToGroupCommit.GroupId)
	if err != nil {
		tx.Rollback()
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}

	_, err = tx.Exec("update `group` set people_num=people_num+1 where id=?", addToGroupCommit.GroupId)
	if err != nil {
		tx.Rollback()
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	tx.Commit()
	common.ReturnResponse(c, 200, 200, "success", nil)
}
