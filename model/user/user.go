package user

type User struct {
	Uid      int    `gorm:"type:primary_key;AUTO_INCREMENT" json:"uid" db:"uid"`
	Rname    string `gorm:"type:varchar(100);not null;" json:"rname" db:"rname"`
	Username string `gorm:"type:varchar(100);not null;" json:"username" db:"username"`
	Mobile   string `gorm:"type:varchar(100);not null;" json:"mobile" db:"mobile"`
}

type UserFriend struct {
	IsDel       uint8  `db:"is_del" json:"is_del"`
	Uid         int    `json:"uid" db:"uid"`
	Id          int    `json:"id" db:"id"`
	FriendUid   int    `db:"friend_uid" json:"friend_uid"`
	CreatedTime uint64 `db:"created_time" json:"created_time"`
	UpdateTime  uint64 `db:"created_time" json:"created_time"`
}

type UserAddFriendRequest struct {
	Id          int    `json:"id" db:"id"`
	RequestUid  int    `json:"request_uid" db:"request_uid"`
	ReceiveUid  int    `json:"receive_uid" db:"receive_uid"`
	CreatedTime uint64 `db:"created_time" json:"created_time"`
}
