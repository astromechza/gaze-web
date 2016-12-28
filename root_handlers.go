package main

import (
	"github.com/kataras/iris"
)

func indexHandler(ctx *iris.Context) {
	ctx.MustRender("root/index.html", struct{ Title string }{"Home"})
}
