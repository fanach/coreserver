package api

import (
	"fmt"

	"gopkg.in/kataras/iris.v6"
)

// GetRoot handles GET /
func GetRoot(ctx *iris.Context) {
	fmt.Println(ctx.Session().GetString(KeySessID))
	fmt.Println(ctx.Session().GetString(KeySessUsername))
	ctx.WriteString("Fanach core server")
}
