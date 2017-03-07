package main

import (
	"strconv"
	"time"

	"gopkg.in/kataras/iris.v6"
)

type loggerMiddleware struct {
}

func (m loggerMiddleware) Serve(ctx *iris.Context) {
	startTime := time.Now()
	ctx.Next()
	latency := time.Now().Sub(startTime).Seconds() * 1000
	ip := ctx.RemoteAddr()
	status := strconv.Itoa(ctx.ResponseWriter.StatusCode())
	path := ctx.Path()
	qs := string(ctx.Request.URL.RawQuery)
	method := ctx.Method()

	ctx.Log(iris.DevMode, "[%v - %.2fms] %s %s %s?%s \n", status, latency, ip, method, path, qs)
	ctx.Next()
}
