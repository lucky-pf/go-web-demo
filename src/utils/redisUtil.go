package utils

import (
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"time"
)

var Pool *redigo.Pool

func init() {
	Pool = PoolInitRedis("127.0.0.1:6379", "")
	if Pool != nil {
		fmt.Println("REDIS连接成功...")
	} else {
		fmt.Println("REDIS连接失败...")
	}
}

// PoolInitRedis redis pool
func PoolInitRedis(server string, password string) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     2, //空闲数
		IdleTimeout: 240 * time.Second,
		MaxActive:   8, //最大数
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
