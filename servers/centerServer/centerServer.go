package centerServer

//以后改成中心服务器
import (
	"bytes"
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
	long := 2
	buf := make([]byte, long)
	myBuffer := new(bytes.Buffer)
	for {
		cnt, err := conn.Read(buf)
		if cnt == 0 || err != nil {
			_ = conn.Close()
			break
		}
		myBuffer.Write(buf[0:cnt])
		for cnt == long {
			cnt, err = conn.Read(buf)
			if cnt == 0 || err != nil {
				_ = conn.Close()
				break
			}
			myBuffer.Write(buf[0:cnt])
		}

		fmt.Println("服务端接收到的命令 : ", strings.TrimSpace(myBuffer.String()))
	}
}
