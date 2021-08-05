package clientServer

//每个节点连接中心服务器用的
import (
	"fmt"
	"net"
)

func ConnectServer() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "192.168.3.36:8086")
	if err != nil {
		fmt.Println("解析ip出错", err)
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("连接tcp出错", err)
		return
	}
	sendToServer(conn)

}

func sendToServer(conn *net.TCPConn) {
	word := "Test Send To Server"
	numWrite, err := conn.Write([]byte(word))
	if err != nil {
		fmt.Println("写出错", err)
		return
	}
	buffer := make([]byte, 1024)
	numRead, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("接收出错", err)
		return
	}
	fmt.Println(conn.RemoteAddr().String(), "接收到了", string(buffer[:numRead]), numWrite)
}
