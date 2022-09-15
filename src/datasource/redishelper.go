package datasource

import (
	"fmt"
	// redigo "github.com/garyburd/redigo/redis"
	redigo "github.com/gomodule/redigo/redis"
	"time"
)

// RedisGoPool 全局变量 redis连接池
var RedisGoPool *redigo.Pool

// 启动程序时，初始化连接池
func init() {
	RedisGoPool = initRedisPool()
	if RedisGoPool != nil {
		fmt.Println("REDIS连接池初始化完成...")
	} else {
		fmt.Println("REDIS连接池初始化失败...")
	}
}

func initRedisPool() (pool *redigo.Pool) {
	server := "127.0.0.1:6379"
	password := ""
	db := 0
	pool = &redigo.Pool{
		MaxIdle:     8,                 // 最大空闲连接数
		MaxActive:   0,                 // 表示和数据库的最大连接，0表示没有限制
		IdleTimeout: 240 * time.Second, // 最大空闲时间
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
			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
	return
}

// SetExpTime 设置过期时间，key: key expTime: 过期时间，秒
func SetExpTime(key string, expTime int) {
	conn := RedisGoPool.Get()
	defer conn.Close()
	_, err := conn.Do("expire", key, expTime)
	if err != nil {
		fmt.Println("set expire error: ", err)
		return
	}
}

func LPush(key string, value ...string) {
	c := RedisGoPool.Get()
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
	c := RedisGoPool.Get()
	defer c.Close()
	_, err := c.Do("del", key)
	if err != nil {
		fmt.Println("del error: ", err)
		return
	}
}

func LPushList(key string, values []string) {
	c := RedisGoPool.Get()
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
	c := RedisGoPool.Get()
	defer c.Close()
	values, err := redigo.Values(c.Do("lrange", key, 0, -1))
	if err != nil {
		fmt.Println("ltrim error: ", err)
		return
	}
	for _, v := range values {
		// fmt.Printf("%s ", v.([]byte))
		// rs = append(rs, v.(string))
		strV, _ := redigo.String(v, nil)
		rs = append(rs, strV)
	}
	return
}

func LPop(key string) string {
	c := RedisGoPool.Get()
	defer c.Close()
	r, err := redigo.String(c.Do("lpop", key))
	if err != nil {
		fmt.Println("lpush error: ", err)
		return ""
	}
	return r
}
