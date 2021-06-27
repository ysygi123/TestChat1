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
func AddGroup(c *gin.Context) {
	groupCreateValidate := &groupvalidate.GroupCreateValidate{}
	if err := common.AutoValidate(c, groupCreateValidate); err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	t := uint64(time.Now().Unix())
	result, err := mysql.DB.Exec("INSERT INTO `group` (`group_name`,`created_uid`,`created_time`,`update_time`) VALUES (?,?,?,?)",
		groupCreateValidate.GroupName, groupCreateValidate.CreatedUid, t, t)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	groupId, err := result.LastInsertId()
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
	_, err = mysql.DB.Exec("INSERT INTO `group_users` (`uid`,`group_id`, `created_time`,`update_time`) VALUES (?,?,?,?)",
		groupCreateValidate.CreatedUid, groupId, t, t)
	if err != nil {
		common.ReturnResponse(c, 200, 400, err.Error(), nil)
		return
	}
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