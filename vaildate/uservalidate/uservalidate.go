package uservalidate

type LoginValidate struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func (r *LoginValidate)GetError()  {
	
}