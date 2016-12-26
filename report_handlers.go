package main

import (
	"github.com/kataras/iris"
)

func newReportHandler(ctx *iris.Context) {
	ctx.EmitError(iris.StatusInternalServerError)
}

func listReportsHandler(ctx *iris.Context) {
	ctx.EmitError(iris.StatusInternalServerError)
}

func getReportHandler(ctx *iris.Context) {
	ctx.EmitError(iris.StatusInternalServerError)
}
