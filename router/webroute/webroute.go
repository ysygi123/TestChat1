package webroute

import (
	"TestChat1/controller/webcontroller"
	"TestChat1/servers/web"
)

func SetWebRoute() {
	web.GinEniger.POST("/user/Login", webcontroller.Login)
	web.GinEniger.GET("/user/LookClient", webcontroller.LookClient)
	web.GinEniger.POST("/user/AuthClient", webcontroller.AuthClient)
	web.GinEniger.POST("/user/SendMessage", webcontroller.SendMessage)
}
