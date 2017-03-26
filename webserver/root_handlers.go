package webserver

import "gopkg.in/kataras/iris.v6"

func IndexHandler(ctx *iris.Context) {
	ctx.MustRender("root/index.html", struct {
		Title string
	}{"Home"})
}
