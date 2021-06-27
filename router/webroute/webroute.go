package webroute

import (
	"TestChat1/controller/webcontroller"
	"TestChat1/middleware"
	"TestChat1/servers/web"
)

func SetWebRoute() {
	web.GinEniger.POST("/user/Login", webcontroller.Login)
	web.GinEniger.Use(middleware.AuthSession())
	//这样写会访问两次这个中间件。。很奇怪
	//web.GinEniger.Use(middleware.AuthSession()).GET("/user/LookClient", webcontroller.LookClient)
	web.GinEniger.GET("/user/LookClient", webcontroller.LookClient)
	web.GinEniger.POST("/user/AuthClient", webcontroller.AuthClient)
	web.GinEniger.POST("/user/AddFriendCommit", webcontroller.AddFriendCommit)
	web.GinEniger.GET("/user/GetFriendsList/:uid", webcontroller.GetFriendsList)
	web.GinEniger.GET("/message/GetMessageList/:uid", webcontroller.GetMessageList)
	web.GinEniger.POST("/message/SendMessage", webcontroller.SendMessage)
	web.GinEniger.POST("/group/CreateGroup", webcontroller.CreateGroup)
	web.GinEniger.GET("/group/GroupList/:uid", webcontroller.GroupList)
	web.GinEniger.POST("/group/AddToGroupCommit", webcontroller.AddToGroupCommit)
}
