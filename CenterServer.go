package main

import (
	"TestChat1/servers/configStart"
	"TestChat1/servers/serverChat/centerServer"
)

// 想尝试写一个tcp中心服务器，然后各个子服务器发消息的时候通过中心服务器来判断
// 这些人在哪些子服务器上面
func main() {
	configStart.ConfigStart()
	centerServer.CenterServerStart()
}
