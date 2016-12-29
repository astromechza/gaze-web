package main

import (
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"

	"strconv"

	"github.com/kataras/iris"
)

func paramIntOrDefault(ctx *iris.Context, name string, def int64) int64 {
	i, e := ctx.URLParamInt(name)
	if e == nil {
		return int64(i)
	}
	return def
}

func paginationReadyQueryString(q *url.Values) string {
	s := q.Encode()
	if len(s) > 0 {
		s = s + "&"
	}
	return "?" + s
}

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
	pageNum := paramIntOrDefault(ctx, "page", 1)
	numberPerPage := paramIntOrDefault(ctx, "numberperpage", 25)

	if pageNum < 1 {
		pageNum = 1
	}
	if numberPerPage < 5 {
		numberPerPage = 5
	}
	if numberPerPage > 500 {
		numberPerPage = 500
	}

	hostname := strings.TrimSpace(ctx.URLParam("hostname"))
	cmdName := strings.TrimSpace(ctx.URLParam("name"))
	exitType := strings.TrimSpace(ctx.URLParam("exit"))
	exitType = strings.ToLower(exitType)
	exitCode, err := strconv.ParseInt(exitType, 10, 64)
	if err != nil {
		if exitType == "zero" || exitType == "success" {
			exitCode = 0
		} else if exitType == "nonzero" || exitType == "failure" {
			exitCode = -1
		}
	}

	var reports []Report
	var totalRecords int64
	q := Database.Active.Order("ulid desc")
	q = q.Offset((pageNum - 1) * numberPerPage).Limit(numberPerPage)
	if hostname != "" {
		q = q.Where("hostname = ?", hostname)
	}
	if cmdName != "" {
		q = q.Where("name = ?", cmdName)
	}
	if exitType != "" {
		if exitCode < 0 {
			q = q.Where("exit_code > 0")
		} else {
			q = q.Where("exit_code = ?", exitCode)
		}
	}
	q.Find(&reports)
	q = Database.Active.Table("reports")
	if hostname != "" {
		q = q.Where("hostname = ?", hostname)
	}
	if cmdName != "" {
		q = q.Where("name = ?", cmdName)
	}
	if exitType != "" {
		if exitCode < 0 {
			q = q.Where("exit_code > 0")
		} else {
			q = q.Where("exit_code = ?", exitCode)
		}
	}
	q.Count(&totalRecords)
	var tags []Tag
	Database.Active.Find(&tags)

	activeQuery := make(url.Values)
	if hostname != "" {
		activeQuery.Add("hostname", hostname)
	}
	if cmdName != "" {
		activeQuery.Add("name", cmdName)
	}
	if exitType != "" {
		activeQuery.Add("exit", exitType)
	}

	numberOfPages := int64(math.Ceil(float64(totalRecords) / float64(numberPerPage)))

	ctx.MustRender("reports/list.html", struct {
		Title         string
		Reports       []Report
		Tags          []Tag
		TotalRecords  int64
		CurrentPage   int64
		TotalPages    int64
		UrlToPaginate string
	}{"Reports", reports, tags, totalRecords, pageNum, numberOfPages, "/reports" + paginationReadyQueryString(&activeQuery)})
}

func getReportHandler(ctx *iris.Context) {
	u := ctx.Param("ulid")
	var report Report
	if err := Database.Active.Preload("Tags").Where("ulid = ?", u).First(&report).Error; err != nil {
		ctx.EmitError(iris.StatusNotFound)
		return
	}

	ctx.MustRender("reports/show.html", struct {
		Title         string
		Report        Report
		ElapsedString string
	}{report.Ulid, report, fmtElapsedTime(report.ElapsedSeconds)})
}
