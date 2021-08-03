package centerServer

//以后改成中心服务器
import (
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
	buf := make([]byte, 4096)
	for {
		cnt, err := conn.Read(buf)
		if cnt == 0 || err != nil {
			_ = conn.Close()
			break
		}
		fmt.Println("服务端接收到的命令 : ", strings.TrimSpace(string(buf[0:cnt])))
	}
}
