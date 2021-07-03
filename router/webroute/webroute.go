package webroute

import (
	"TestChat1/controller/webcontroller"
	"TestChat1/middleware"
	"TestChat1/servers/web"
)

func SetWebRoute() {
	web.GinEniger.POST("/user/Login", webcontroller.Login)
	web.GinEniger.GET("/user/LookClient", webcontroller.LookClient)
	web.GinEniger.Use(middleware.AuthSession())
	//这样写会访问两次这个中间件。。很奇怪
	//web.GinEniger.Use(middleware.AuthSession()).GET("/user/LookClient", webcontroller.LookClient)
	//web.GinEniger.POST("/user/AuthClient", webcontroller.AuthClient)
	web.GinEniger.POST("/user/AddFriendCommit", webcontroller.AddFriendCommit)      //同意好友请求
	web.GinEniger.GET("/user/GetFriendsList/:uid", webcontroller.GetFriendsList)    //获取好友列表
	web.GinEniger.GET("/message/GetMessageList/:uid", webcontroller.GetMessageList) //获取消息列表
	web.GinEniger.GET("/message/GetSelfChat", webcontroller.GetSelfChat)            //私聊消息
	web.GinEniger.GET("/message/GetGroupChat", webcontroller.GetGroupChat)          //群聊消息
	web.GinEniger.POST("/message/SendMessage", webcontroller.SendMessage)           //发送所有消息
	web.GinEniger.POST("/group/CreateGroup", webcontroller.CreateGroup)             //创建群
	web.GinEniger.GET("/group/GroupList/:uid", webcontroller.GroupList)             //查看群列表
	web.GinEniger.POST("/group/AddToGroupCommit", webcontroller.AddToGroupCommit)   //同意加群
}
