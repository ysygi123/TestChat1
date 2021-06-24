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

func (r *LoginValidate) GetError() {

}
