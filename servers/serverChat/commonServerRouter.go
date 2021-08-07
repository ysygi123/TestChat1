package serverChat

import "net"

//客户和主服务器之间相互交流

type ClientRequestMsg struct {
	Cmd    string                 `json:"cmd"`
	Params map[string]interface{} `json:"params"`
}

type ServerResponseMsg struct {
	Cmd    string                 `json:"cmd"`
	Params map[string]interface{} `json:"params"`
}

var ServerRouteManager *ServerRoute

type ServerRouteFunction func(smsg *ClientRequestMsg, conn net.Conn)

type ServerRoute struct {
	Route map[string]*ServerRouteFunction
}

func NewServerRouteManager() {
	ServerRouteManager = &ServerRoute{}
}

//注册路由
func (this *ServerRoute) Register(commandName string, handler *ServerRouteFunction) {
	this.Route[commandName] = handler
}

//获取路由函数
func (this *ServerRoute) GetHandler(commandName string) *ServerRouteFunction {
	h, ok := this.Route[commandName]
	if !ok {
		return nil
	}
	return h
}
