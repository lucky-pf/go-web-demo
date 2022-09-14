package main

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"go-web-demo/src/service"
)

func main() {
	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("<b>Hello!</b>")
	})

	app.Get("/users", func(ctx iris.Context) {
		users := service.GetUserList()
		for _, v := range users {
			fmt.Println(v)
		}
		jsons, errs := json.Marshal(users)
		if errs != nil {
			fmt.Println(errs.Error())
		}

		ctx.HTML("<b>" + string(jsons) + "</b>")
	})

	app.Get("/set", func(ctx iris.Context) {
		name := ctx.URLParam("name")
		service.SetStr(name)
		name = service.GetStr()
		ctx.HTML("<b> " + name + "!</b>")
	})

	app.Run(iris.Addr(":8899"), iris.WithConfiguration(iris.YAML("./src/iris.yml")))
}
