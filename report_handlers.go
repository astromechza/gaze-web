package main

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/url"
	"strings"
	"time"

	"strconv"

	"github.com/oklog/ulid"
	"gopkg.in/kataras/iris.v6"
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

func failBadRequest(ctx *iris.Context, err error) {
	ctx.Log(iris.DevMode, "Bad request due to: %v", err.Error())
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Log(iris.DevMode, "Data was: \n%v", string(body))
	ctx.EmitError(400)
	ctx.WriteString("Bad Request")
}

type jsonErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func failJSON(ctx *iris.Context, code int, reason string) {
	ctx.Log(iris.DevMode, "Failure (%v) due to: '%v'", code, reason)
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Log(iris.DevMode, "Data was: \n%v", string(body))
	ctx.EmitError(code)
	s, _ := json.Marshal(jsonErrorMessage{code, reason})
	ctx.Write(s)
}

func newReportHandler(ctx *iris.Context) {
	var report Report
	err := ctx.ReadJSON(&report)
	if err != nil {
		failJSON(ctx, 400, "Invalid json payload")
		return
	}

	nowTime := time.Now()

	// fields you can't provide
	report.ID = 0
	report.CreatedAt = nowTime
	report.UpdatedAt = nowTime
	report.DeletedAt = nil

	// fields you must provide
	report.Hostname = strings.TrimSpace(strings.ToLower(report.Hostname))
	if report.Hostname == "" {
		failJSON(ctx, 400, "Missing 'hostname' value")
		return
	}
	report.Name = strings.TrimSpace(report.Name)
	if report.Name == "" {
		failJSON(ctx, 400, "Missing 'name' value")
		return
	}

	// fields that are optional
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
						ctx.Log(iris.DevMode, "db create tag error %v", err.Error())
						ctx.EmitError(400)
						return
					}
				}
				report.Tags = append(report.Tags, tag)
			}
		}
	}
	if report.ElapsedSeconds < 0 {
		failJSON(ctx, 400, "'elapsed_seconds' value must be >= 0")
		return
	}
	if report.EndTime.Before(report.StartTime) {
		failJSON(ctx, 400, "'end_time' value must be >= 'start_time' value")
		return
	}
	if report.StartTime.IsZero() {
		if report.EndTime.IsZero() {
			report.StartTime = nowTime
		} else {
			report.StartTime = report.EndTime.Add(time.Duration(int64(-report.ElapsedSeconds * float32(time.Second))))
		}
	}
	if report.EndTime.IsZero() {
		report.EndTime = report.StartTime.Add(time.Duration(int64(-report.ElapsedSeconds * float32(time.Second))))
	}
	if report.Ulid == "" {
		report.Ulid = ulid.MustNew(ulid.Timestamp(time.Now()), RandomSource).String()
	}

	if err = Database.Active.Create(&report).Error; err != nil {
		ctx.Log(iris.DevMode, "db create report error "+err.Error())
		failJSON(ctx, 500, "Server Error")
	} else {
		ctx.SetStatusCode(204)
	}
}

func listReportsHandler(ctx *iris.Context) {
	pageNum := paramIntOrDefault(ctx, "page", 1)
	numberPerPage := paramIntOrDefault(ctx, "numberperpage", 50)

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
		} else if exitType == "any" || exitType == "anything" {
			exitType = ""
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

	var lastFailure Report
	q = Database.Active.Order("ulid desc")
	if hostname != "" {
		q = q.Where("hostname = ?", hostname)
	}
	if cmdName != "" {
		q = q.Where("name = ?", cmdName)
	}
	q = q.Where("exit_code > ?", 0)
	q.First(&lastFailure)

	var lastSuccess Report
	q = Database.Active.Order("ulid desc")
	if hostname != "" {
		q = q.Where("hostname = ?", hostname)
	}
	if cmdName != "" {
		q = q.Where("name = ?", cmdName)
	}
	q = q.Where("exit_code = 0")
	q.First(&lastSuccess)

	numberOfPages := int64(math.Ceil(float64(totalRecords) / float64(numberPerPage)))

	ctx.MustRender("reports/list.html", struct {
		Title        string
		Reports      []Report
		LastSuccess  Report
		LastFailure  Report
		Tags         []Tag
		TotalRecords int64

		CurrentPage   int64
		TotalPages    int64
		URLToPaginate string

		FormName     string
		FormHostname string
		FormExit     string
	}{
		"Reports",
		reports, lastSuccess, lastFailure,
		tags,
		totalRecords, pageNum, numberOfPages,
		"/reports" + paginationReadyQueryString(&activeQuery),
		cmdName,
		hostname,
		exitType,
	})
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
