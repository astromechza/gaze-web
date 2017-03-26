package webserver

import (
	"path/filepath"

	"github.com/AstromechZA/gaze-web/utils"

	iris "gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/view"
)

func Setup(srvDir string) *iris.Framework {
	app := iris.New()

	app.Adapt(iris.DevLogger())
	app.Adapt(httprouter.New())

	engine := view.HTML(filepath.Join(srvDir, "templates"), ".html")
	engine.Layout("root/layout.html")
	engine.Funcs(utils.BuildTemplateFuncsMap())
	app.Adapt(engine)

	app.Use(LoggerMiddleware{})

	app.StaticWeb("/static", filepath.Join(srvDir, "static"))

	app.Get("/", IndexHandler)
	app.Post("/report", NewReportHandler)
	app.Put("/report", NewReportHandler)
	app.Get("/reports", ListReportsHandler)
	app.Get("/reports/:ulid", GetReportHandler)
	app.Get("/graph", GraphReportsHandler)

	app.OnError(iris.StatusInternalServerError, Error500Handler)
	return app
}
