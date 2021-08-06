package clientServer

//每个节点连接中心服务器用的
import (
	"TestChat1/servers/serverChat"
	"fmt"
	"net"
	"strings"
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
	word := "Test Send To Server"
	_, err := conn.Write([]byte(word))
	if err != nil {
		fmt.Println("写出错", err)
		return
	}

	for {
		iobuffer, err := serverChat.ConnReadMany(conn)
		if err != nil {
			fmt.Println("读出错", err)
			return
		}
		fmt.Println("客户端接收到的命令 : ", strings.TrimSpace(iobuffer.String()))
	}
}
