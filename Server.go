package main

import (
	"TestChat1/controller/websocketcontroller"
	"TestChat1/db/mysql"
	"TestChat1/db/redis"
	"TestChat1/router/webroute"
	"TestChat1/servers/backtask"
	"TestChat1/servers/idDispatch"
	"TestChat1/servers/web"
	"TestChat1/servers/websocket"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
)

func main() {
	initConfig()
	mysql.NewMysqlDB()
	redis.NewRedisCluster()
	websocket.ClientMangerInstanceInit()
	websocket.WebSocketRouteMangerInit()
	websocket.WebSocketRouteManger.AllRegisterRoute()
	go websocket.ClientMangerInstance.LoopToKillChild()
	go backtask.AllBackTask()
	err := idDispatch.NewWorker(1)
	if err != nil {
		panic(err)
	}
	go func() {
		http.HandleFunc("/ws", websocketcontroller.FirstPage)
		http.ListenAndServe(":8087", nil)
	}()
	webroute.SetWebRoute()
	_ = web.GinEniger.Run(":8088")
}

func initConfig() {
	viper.SetConfigName("test")
	viper.AddConfigPath("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	fmt.Println(viper.GetString("mysql.username"))
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
