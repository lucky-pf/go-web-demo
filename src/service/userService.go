package service

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go-web-demo/src/utils"
)

func SetStr(str string) {
	coon := utils.Pool.Get()
	defer coon.Close()

	_, err := coon.Do("SET", "name", str)
	if err != nil {
		fmt.Println("redis 设值失败:", err)
	}

}

func GetStr() (name string) {
	coon := utils.Pool.Get()
	defer coon.Close()

	// rs, err := redis.Strings(coon.Do("GET", "name"))
	name, err := redis.String(coon.Do("GET", "name"))
	// rs, err := coon.Do("GET", "name")
	if err != nil {
		fmt.Println("redis取值失败：", err)
		return
	}
	fmt.Println("get key: ", name)
	return
}
