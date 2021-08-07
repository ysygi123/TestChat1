package centerService

import (
	"TestChat1/servers/serverChat"
	"fmt"
	"net"
)

func Test1(smsg *serverChat.ClientRequestMsg, conn net.Conn) {
	fmt.Println(smsg.Cmd, smsg.Params)
	conn.Write([]byte("testt11111"))
}
