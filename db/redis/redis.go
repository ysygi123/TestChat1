package redis

import (
	redgio "github.com/garyburd/redigo/redis"
)

var RedisPool *redgio.Pool

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
	}

}
