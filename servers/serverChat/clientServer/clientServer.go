package clientServer

//每个节点连接中心服务器用的
import (
	"TestChat1/servers/serverChat"
	"encoding/json"
	"fmt"
	"net"
)

func ConnectServer() {
	conn, err := net.Dial("tcp", "192.168.3.36:8086")
	if err != nil {
		fmt.Println("解析ip出错", err)
		return
	}
	sendToServer(conn)
}

func sendToServer(conn net.Conn) {
	s := serverChat.ClientRequestMsg{
		Cmd: "test1",
		Params: map[string]interface{}{
			"a1": "去你妈的",
		},
	}
	requestByte, err := json.Marshal(s)
	_, err = conn.Write(requestByte)
	if err != nil {
		fmt.Println("写出错", err)
		return
	}

	for {
		returnByte, err := serverChat.ConnReadMany(conn)
		if err != nil {
			fmt.Println("客户读出错", err)
			return
		}
		fmt.Println("客户端接收到的命令 : ", string(returnByte))
	}
}
