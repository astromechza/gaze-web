package webserver

import (
	"github.com/AstromechZA/gaze-web/random"

	"gopkg.in/kataras/iris.v6"
)

func logThatAnErrorOccured(code int, message string, ctx *iris.Context) string {
	u := random.NewUlidNow().String()
	ctx.Log(iris.DevMode, "Error %v (ulid: %v) during %v %v: %v", code, u, string(ctx.Method()), string(ctx.Path()), message)
	return u
}

func Error500Handler(ctx *iris.Context) {
	u := logThatAnErrorOccured(500, "not implemented", ctx)
	ctx.MustRender("root/500.html", struct{ Ulid string }{u})
}
