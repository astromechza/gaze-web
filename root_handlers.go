package main

import "gopkg.in/kataras/iris.v6"

func indexHandler(ctx *iris.Context) {
	ctx.MustRender("root/index.html", struct {
		Title string
	}{"Home"})
}
