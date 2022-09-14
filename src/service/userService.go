package service

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go-web-demo/src/datasource"
	"go-web-demo/src/model"
)

func SetStr(str string) {
	coon := datasource.REDIGO_POOL.Get()
	defer coon.Close()

	_, err := coon.Do("SET", "name", str)
	if err != nil {
		fmt.Println("redis 设值失败:", err)
	}

}

func GetStr() (name string) {
	coon := datasource.REDIGO_POOL.Get()
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

func GetUserList() []model.User {
	db := datasource.GetMySqlDB()
	var users []model.User
	db.Find(&users)
	// fmt.Printf("user:%#v\n", users)
	return users
}
