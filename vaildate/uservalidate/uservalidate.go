package uservalidate

type LoginValidate struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type Auth struct {
	Uid     int    `from:"uid" binding:"required"`
	Session string `from:"session" binding:"required"`
}

func (r *LoginValidate) GetError() {

}
