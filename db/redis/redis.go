package redis

import (
	redgio "github.com/garyburd/redigo/redis"
)

var RedisPool *redgio.Pool

//连接过多的时候会出现 redigo: connection pool exhausted 需要及时close 好像 GET方法里面也有些 must close
func NewRedisDB() {
	RedisPool = &redgio.Pool{
		MaxIdle:   1000,
		MaxActive: 100,
		Dial: func() (conn redgio.Conn, e error) {
			conn, e = redgio.Dial("tcp", "127.0.0.1:6379")
			if e != nil {
				return nil, e
			}
			return
		},
		Wait: true,
	}

}
