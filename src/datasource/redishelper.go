package datasource

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

// POOL 全局变量 redis连接池
var POOL *redis.Pool

// 启动程序时，初始化连接池
func init() {
	server := "202.46.45.219:26379"
	password := "sjdnjka2123"
	db := 0
	POOL = &redis.Pool{
		MaxIdle:     8,                 // 最大空闲连接数
		MaxActive:   0,                 // 表示和数据库的最大连接，0表示没有限制
		IdleTimeout: 240 * time.Second, // 最大空闲时间
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
}

// SetExpTime 设置过期时间，key: key expTime: 过期时间，秒
func SetExpTime(key string, expTime int) {
	conn := POOL.Get()
	_, err := conn.Do("expire", key, expTime)
	if err != nil {
		fmt.Println("set expire error: ", err)
		return
	}
}

func LPush(key string, value ...string) {
	c := POOL.Get()
	defer c.Close()
	_, err := c.Do("lpush", key, value)
	if err != nil {
		fmt.Println("lpush error: ", err)
		return
	}
}

/*func LPushList(key string, values []string) {
	c := POOL.Get()
	defer c.Close()
	_, err := c.Do("lpush", key, values)
	if err != nil {
		fmt.Println("lpush error: ", err)
		return
	}
}*/

func DelKey(key string) {
	c := POOL.Get()
	defer c.Close()
	_, err := c.Do("del", key)
	if err != nil {
		fmt.Println("del error: ", err)
		return
	}
}

func LPushList(key string, values []string) {
	c := POOL.Get()
	defer c.Close()

	for _, v := range values {
		// 把命令写到输出缓冲区
		c.Send("lpush", key, v)
	}
	// 把缓冲区的命令刷新到redis服务器
	c.Flush()
	// 接收redis给予的响应
	_, err := c.Receive()
	if err != nil {
		fmt.Println("lpush error: ", err)
		return
	}
}

func GetAllList(key string) (rs []string) {
	c := POOL.Get()
	defer c.Close()
	values, err := redis.Values(c.Do("lrange", key, 0, -1))
	if err != nil {
		fmt.Println("ltrim error: ", err)
		return
	}
	for _, v := range values {
		// fmt.Printf("%s ", v.([]byte))
		// rs = append(rs, v.(string))
		strV, _ := redis.String(v, nil)
		rs = append(rs, strV)
	}
	return
}

func LPop(key string) string {
	c := POOL.Get()
	defer c.Close()
	r, err := redis.String(c.Do("lpop", key))
	if err != nil {
		fmt.Println("lpush error: ", err)
		return ""
	}
	return r
}
