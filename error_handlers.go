package main

import (
	"time"

	"github.com/oklog/ulid"
	"gopkg.in/kataras/iris.v6"
)

func logThatAnErrorOccured(code int, message string, ctx *iris.Context) string {
	u := ulid.MustNew(ulid.Timestamp(time.Now()), RandomSource).String()
	ctx.Log(iris.DevMode, "Error %v (ulid: %v) during %v %v: %v", code, u, string(ctx.Method()), string(ctx.Path()), message)
	return u
}

func error500Handler(ctx *iris.Context) {
	u := logThatAnErrorOccured(500, "not implemented", ctx)
	ctx.MustRender("root/500.html", struct{ Ulid string }{u})
}
