package group

type Group struct {
	IsDel       uint8  `db:"is_del" json:"is_del"`
	PeopleNum   int16  `db:"people_num" json:"people_num"`
	Id          int    `db:"id" json:"id"`
	CreatedUid  int    `db:"created_uid" json:"created_uid"`
	CreatedTime uint64 `db:"created_time" json:"created_time"`
	UpdateTime  uint64 `db:"created_time" json:"created_time"`
	GroupName   string `db:"group_name" json:"group_name"`
}

type GroupMessage struct {
	Id             int    `db:"id" json:"id"`
	GroupId        int    `db:"group_id" json:"group_id"`
	CreatedTime    uint64 `db:"created_time" json:"created_time"`
	MessageContent string `db:"message_content" json:"message_content"`
}

type GroupUsers struct {
	IsDel       int8 `db:"is_del" json:"is_del"` // 1不删2删
	Id          int  `db:"id" json:"id"`
	CreatedTime int  `db:"created_time" json:"created_time"`
	UpdateTime  int  `db:"update_time" json:"update_time"`
	Uid         int  `db:"uid" json:"uid"`
	GroupId     int  `db:"group_id" json:"group_id"`
}
