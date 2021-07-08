package main

import (
	"TestChat1/db/redis"
	"fmt"
	tt "github.com/go-redis/redis"
	"strconv"
)

func main() {
	brpopTest()
}

func testRedis() {
	redisdb := tt.NewClient(&tt.Options{
		Addr:     "localhost:6379", // use default Addr
		Password: "",               // no password set
		DB:       0,                // use default DB
	})
	a, _ := redisdb.HGet("468afced1737079cb46b2c113ec0a6841", "uid").Result()
	fmt.Println(a == "")
}

func clusterTest() {
	pong, err := redis.GoRedisCluster.Ping().Result()
	fmt.Println(pong, err)
	fmt.Println("pool state init state:", redis.GoRedisCluster.PoolStats())
	for i := 0; i < 10; i++ {
		redis.GoRedisCluster.Set(strconv.Itoa(i), i, 0)
	}
	redis.GoRedisCluster.Close()
}

func brpopTest() {
	for {
		res, err := redis.GoRedisCluster.BRPop(0, "test_queue").Result()
		fmt.Println(res, err)
	}
}
