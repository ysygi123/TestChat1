package main

import (
	"TestChat1/db/mysql"
	"TestChat1/db/redis"
	"fmt"
	"strconv"
)

func main() {
	mysql.NewMysqlDB()
	redis.NewRedisDB()
	allUid := make([]int, 0)
	for i := 0; i < 10; i++ {
		table := fmt.Sprintf("user_login_%d", i)
		rows, _ := mysql.DB.Query("select uid from " + table)
		tmpUid := 0
		for rows.Next() {
			rows.Scan(&tmpUid)
			allUid = append(allUid, tmpUid)
		}
	}
	for _, v := range allUid {
		session, _ := redis.GoRedisCluster.Get("uidlogin:" + strconv.Itoa(v)).Result()
		fmt.Println(session)
		res, err := redis.GoRedisCluster.Del(session).Result()
		res, err = redis.GoRedisCluster.Del("uidlogin:" + strconv.Itoa(v)).Result()
		fmt.Println(res, err)
	}
}
