package uservalidate

type LoginValidate struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type Auth struct {
	Uid     int    `from:"uid" binding:"required"`
	Session string `from:"session" binding:"required"`
}

type AddFriendRequest struct {
	SendUid    int    `json:"send_uid"`
	ReceiveUid int    `json:"receive_uid"`
	Rname      string `json:"rname"`
}

type AddFriendCommit struct {
	Uid       int `json:"uid"  binding:"required"`
	MessageId int `json:"message_id" binding:"required"`
}

type RegisterValidate struct {
	Username string `json:"username" binding:"required"`
	Passwd   string `json:"passwd" binding:"required"`
}

func (r *LoginValidate) GetError() {

}
