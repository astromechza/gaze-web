package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/kataras/iris"
)

func newReportHandler(ctx *iris.Context) {
	var report Report
	err := ctx.ReadJSON(&report)
	if err != nil {
		fmt.Println(err.Error())
		ctx.EmitError(400)
		return
	}
	report.ID = 0
	report.CreatedAt = time.Now()
	report.UpdatedAt = time.Now()
	report.DeletedAt = nil
	if len(report.RawCommand) > 0 {
		commandS := ""
		for _, part := range report.RawCommand {
			if strings.Contains(part, " ") {
				part = "\"" + part + "\""
			}
			if len(commandS) == 0 {
				commandS += part
			} else {
				commandS += " " + part
			}
		}
		report.Command = commandS
	}
	if len(report.RawTags) > 0 {
		report.Tags = make([]Tag, 0)
		for _, ts := range report.RawTags {
			ts = strings.TrimSpace(strings.ToLower(ts))
			if len(ts) > 0 {
				// validate
				var tag Tag
				if err := Database.Active.Where("text = ?", ts).First(&tag).Error; err != nil {
					tag.Text = ts
					err = Database.Active.Create(&tag).Error
					if err != nil {
						fmt.Println("db create tag error " + err.Error())
						ctx.EmitError(400)
						return
					}
				}
				report.Tags = append(report.Tags, tag)
			}
		}
	}

	if err = Database.Active.Create(&report).Error; err != nil {
		fmt.Println("db create report error " + err.Error())
		ctx.EmitError(400)
	} else {
		ctx.SetStatusCode(204)
	}
}

func listReportsHandler(ctx *iris.Context) {
	var reports []Report
	Database.Active.Table("reports").Find(&reports)
	ctx.MustRender("reports/list.html", struct{ Reports []Report }{reports})
}

func getReportHandler(ctx *iris.Context) {
	u := ctx.Param("ulid")
	var report Report
	if err := Database.Active.Preload("Tags").Where("ulid = ?", u).First(&report).Error; err != nil {
		ctx.EmitError(iris.StatusNotFound)
		return
	}

	ctx.MustRender("reports/show.html", struct {
		Report        Report
		ElapsedString string
	}{report, fmtElapsedTime(report.ElapsedSeconds)})
}
