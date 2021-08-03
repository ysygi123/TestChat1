package main

import (
	"TestChat1/db/mysql"
	"TestChat1/db/redis"
	"TestChat1/servers/userService"
	"TestChat1/test/myrun/testTool"
	"fmt"
	"github.com/spf13/viper"
)

func main() {
	initConfig()
	fmt.Println(viper.Get("mysql"))
	mysql.NewMysqlDB()
	redis.NewRedisCluster()
	c := testTool.GetNewChan(1)
	for i := 0; i < 10; i++ {
		table := fmt.Sprintf("user_login_%d", i)
		rows, _ := mysql.DB.Query("select uid from " + table)
		tmpUid := 0
		for rows.Next() {
			rows.Scan(&tmpUid)
			c.Lock()
			go func() {
				fmt.Println(userService.LoginOut(tmpUid))
				c.Unlock()
			}()
		}
	}
	/*for _, v := range allUid {
		session, _ := redis.GoRedisCluster.Get("uidlogin:" + strconv.Itoa(v)).Result()
		fmt.Println(session)
		res, err := redis.GoRedisCluster.Del(session).Result()
		res, err = redis.GoRedisCluster.Del("uidlogin:" + strconv.Itoa(v)).Result()
		fmt.Println(res, err)
	}*/
}
func initConfig() {
	viper.SetConfigName("test")
	viper.AddConfigPath("../../config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	fmt.Println(viper.GetString("mysql.username"))
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
