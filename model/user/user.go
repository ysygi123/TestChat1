package user

type User struct {
	Uid    uint   `gorm:"type:primary_key;AUTO_INCREMENT" json:"uid" db:"uid"`
	Rname  string `gorm:"type:varchar(100);not null;" json:"rname" db:"rname"`
	Passwd string `gorm:"type:varchar(100);not null;" json:"passwd" db:"passwd"`
}
