package webroute

import (
	"TestChat1/controller/webcontroller"
	"TestChat1/servers/web"
)

func SetWebRoute() {
	web.GinEniger.POST("/user/Login", webcontroller.Login)
	web.GinEniger.POST("/user/LookClient", webcontroller.LookClient)
}
