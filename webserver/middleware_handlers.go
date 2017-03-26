package webserver

import (
	"strconv"
	"time"

	"gopkg.in/kataras/iris.v6"
)

type LoggerMiddleware struct {
}

func (m LoggerMiddleware) Serve(ctx *iris.Context) {
	startTime := time.Now()
	ctx.Next()
	latency := time.Now().Sub(startTime).Seconds() * 1000
	ip := ctx.RemoteAddr()
	status := strconv.Itoa(ctx.ResponseWriter.StatusCode())
	path := ctx.Path()
	qs := string(ctx.Request.URL.RawQuery)
	method := ctx.Method()

	ctx.Log(iris.DevMode, "(%.2fms) %s %s %s?%s [%s]\n", latency, ip, method, path, qs, status)
	ctx.Next()
}
