package centerServer

//以后改成中心服务器
import (
	"TestChat1/servers/serverChat"
	"encoding/json"
	"fmt"
	"net"
)

func CenterServerStart() {
	server, err := net.Listen("tcp", "192.168.3.36:8086")
	fmt.Println(server, err)

	if err != nil {
		fmt.Println("启动失败了啊", err)
		return
	}

	fmt.Println("我试一试开起来看看")
	RegisterCenterServerRoute()
	for {
		conn, err := server.Accept()

		if err != nil {
			fmt.Println("获取conn的时候出错", err)
			continue
		}

		go connHandle(conn)
	}
}

func connHandle(conn net.Conn) {
	if conn == nil {
		fmt.Println("这个conn有问题")
		return
	}

	for {
		returnByte, err := serverChat.ConnReadMany(conn)
		if err != nil {
			fmt.Println("服务端读出错", err)
			return
		}
		s := serverChat.ClientRequestMsg{}

		err = json.Unmarshal(returnByte, &s)

		fmt.Println(s, err)

		if err != nil {
			EchoError("查无此路由"+err.Error(), conn)
			continue
		}
		handle := serverChat.ServerRouteManager.GetServerHandler(s.Cmd)
		fmt.Println(s.Cmd, serverChat.ServerRouteManager.ServerRoute)
		if handle == nil {
			EchoError("查无此路由"+err.Error(), conn)
			continue
		}
		(handle)(&s, conn)
	}
}

func EchoError(errMsg string, conn net.Conn) error {
	msg := serverChat.ServerResponseMsg{
		Cmd: "ERR",
		Params: map[string]interface{}{
			"errMsg": errMsg,
		},
	}

	msgByte, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = conn.Write(msgByte)
	if err != nil {
		return err
	}
	return nil
}
