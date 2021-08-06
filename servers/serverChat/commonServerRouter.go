package serverChat

//客户和主服务器之间相互交流

type ClientRequestMsg struct {
	Cmd string
}

type ServerResponseMst struct {
	Cmd    string
	Params map[string]interface{}
}

var ServerRouteManager *ServerRoute

type ServerRouteFunction func(smsg *map[string]interface{}) *ServerResponseMst

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
