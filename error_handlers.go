package main

import (
	"time"

	"github.com/kataras/iris"
	"github.com/oklog/ulid"
)

func logThatAnErrorOccured(code int, message string, ctx *iris.Context) string {
	u := ulid.MustNew(ulid.Timestamp(time.Now()), RandomSource).String()
	ctx.Log("Error %v (ulid: %v) during %v %v: %v", code, u, string(ctx.Method()), string(ctx.Path()), message)
	return u
}

func error500Handler(ctx *iris.Context) {
	u := logThatAnErrorOccured(500, "not implemented", ctx)
	ctx.MustRender("root/500.html", struct{ Ulid string }{u})
}
