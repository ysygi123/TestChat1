package centerServer

//以后改成中心服务器
import (
	"TestChat1/servers/serverChat"
	"fmt"
	"net"
	"strings"
)

func CenterServerStart() {
	server, err := net.Listen("tcp", "192.168.3.36:8086")
	fmt.Println(server, err)

	if err != nil {
		fmt.Println("启动失败了啊", err)
		return
	}

	fmt.Println("我试一试开起来看看")

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
		iobuffer, err := serverChat.ConnReadMany(conn)
		if err != nil {
			fmt.Println("读出错", err)
			return
		}
		fmt.Println("服务端接收到的命令 : ", strings.TrimSpace(iobuffer.String()))
	}
}
