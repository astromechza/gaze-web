package main

import (
	"strconv"
	"time"

	"github.com/kataras/iris"
)

type loggerMiddleware struct {
}

func (m loggerMiddleware) Serve(ctx *iris.Context) {
	startTime := time.Now()
	ctx.Next()
	latency := time.Now().Sub(startTime).Seconds() * 1000
	ip := ctx.RemoteAddr()
	status := strconv.Itoa(ctx.Response.StatusCode())
	path := ctx.PathString()
	qs := string(ctx.QueryArgs().QueryString())
	method := ctx.MethodString()

	ctx.Log("[%v - %.2fms] %s %s %s?%s \n", status, latency, ip, method, path, qs)

}
