package main

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

func initPool(address string, maxIdle, maxactive int, idleTimeout time.Duration) {
	pool = &redis.Pool{
		MaxIdle:     maxIdle,     //最大空闲链接数
		MaxActive:   maxactive,   //表示和数据库的最大链接数，0表示无限制
		IdleTimeout: idleTimeout, //最大空闲时间，链接之后没有放回，链接会自动增长
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
	}
}
