package centerServer

import (
	"TestChat1/servers/serverChat"
	"TestChat1/servers/serverChat/centerServer/centerService"
)

func RegisterCenterServerRoute() {
	serverChat.ServerRouteManager.RegisterServerRoute("test1", centerService.Test1)
}
