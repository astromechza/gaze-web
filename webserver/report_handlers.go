package webserver

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/url"
	"strings"
	"time"

	"strconv"

	"github.com/AstromechZA/gaze-web/database"
	"github.com/AstromechZA/gaze-web/random"

	"github.com/jinzhu/gorm"
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

func NewReportHandler(ctx *iris.Context) {
	var report database.Report
	err := ctx.ReadJSON(&report)
	if err != nil {
		failJSON(ctx, 400, "Invalid Report payload "+err.Error())
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
		report.Tags = make([]database.Tag, 0)
		for _, ts := range report.RawTags {
			ts = strings.TrimSpace(strings.ToLower(ts))
			if len(ts) > 0 {
				// validate
				var tag database.Tag
				if err := database.ActiveDB.Where("text = ?", ts).First(&tag).Error; err != nil {
					tag.Text = ts
					err = database.ActiveDB.Create(&tag).Error
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
	if report.StartTime.IsZero() {
		if report.EndTime.IsZero() {
			report.StartTime = nowTime
		} else {
			report.StartTime = report.EndTime.Add(time.Duration(int64(-report.ElapsedSeconds * float32(time.Second))))
		}
	}
	if report.EndTime.IsZero() {
		report.EndTime = report.StartTime.Add(time.Duration(int64(report.ElapsedSeconds * float32(time.Second))))
	}

	if report.EndTime.Before(report.StartTime) {
		failJSON(ctx, 400, "'end_time' value must be >= 'start_time' value")
		return
	}

	if report.Ulid == "" {
		report.Ulid = random.NewUlidNow().String()
	}

	if report.ExitDescription == "" {
		report.ExitDescription = "No description provided"
	}

	if err = database.ActiveDB.Create(&report).Error; err != nil {
		ctx.Log(iris.DevMode, "db create report error "+err.Error())
		failJSON(ctx, 500, "Server Error")
	} else {
		ctx.SetStatusCode(204)
	}
}

func parseFilterParams(ctx *iris.Context) (string, string, string, int64) {
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
	return hostname, cmdName, exitType, exitCode
}

func buildURLValues(hostname, cmdName, exitType string) *url.Values {
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
	return &activeQuery
}

func getReportsQuery(hostname, cmdName, exitType string, exitCode int64) *gorm.DB {
	q := database.ActiveDB.Table("reports").Order("start_time desc")
	if hostname != "" {
		q = q.Where("hostname = ?", hostname)
	}
	if cmdName != "" {
		q = q.Where("name = ?", cmdName)
	}
	if exitType != "" {
		q = q.Where("exit_code = ?", exitCode)
	}
	return q
}

func parsePaginationParams(ctx *iris.Context) (int64, int64) {
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
	return pageNum, numberPerPage
}

func ListReportsHandler(ctx *iris.Context) {
	pageNum, numberPerPage := parsePaginationParams(ctx)

	hostname, cmdName, exitType, exitCode := parseFilterParams(ctx)
	activeQuery := buildURLValues(hostname, cmdName, exitType)

	var reports []database.Report
	q := getReportsQuery(hostname, cmdName, exitType, exitCode)
	q = q.Offset((pageNum - 1) * numberPerPage).Limit(numberPerPage)
	q.Find(&reports)

	var totalRecords int64
	getReportsQuery(hostname, cmdName, exitType, exitCode).Count(&totalRecords)

	var lastFailure database.Report
	q = getReportsQuery(hostname, cmdName, exitType, exitCode)
	q = q.Where("exit_code != ?", 0)
	q.First(&lastFailure)

	var lastSuccess database.Report
	q = getReportsQuery(hostname, cmdName, exitType, exitCode)
	q = q.Where("exit_code = 0")
	q.First(&lastSuccess)

	numberOfPages := int64(math.Ceil(float64(totalRecords) / float64(numberPerPage)))

	ctx.MustRender("reports/list.html", struct {
		Title        string
		Reports      []database.Report
		LastSuccess  database.Report
		LastFailure  database.Report
		TotalRecords int64
		GraphLink    string

		CurrentPage   int64
		TotalPages    int64
		URLToPaginate string

		FormName     string
		FormHostname string
		FormExit     string
	}{
		"Reports",
		reports, lastSuccess, lastFailure,
		totalRecords,
		"/graph" + paginationReadyQueryString(activeQuery),
		pageNum, numberOfPages,
		"/reports" + paginationReadyQueryString(activeQuery),
		cmdName,
		hostname,
		exitType,
	})
}

func GetReportHandler(ctx *iris.Context) {
	u := ctx.Param("ulid")
	var report database.Report
	if err := database.ActiveDB.Preload("Tags").Where("ulid = ?", u).First(&report).Error; err != nil {
		ctx.EmitError(iris.StatusNotFound)
		return
	}

	ctx.MustRender("reports/show.html", struct {
		Title  string
		Report database.Report
	}{report.Ulid, report})
}

func GraphReportsHandler(ctx *iris.Context) {
	hostname, cmdName, exitType, exitCode := parseFilterParams(ctx)
	activeQuery := buildURLValues(hostname, cmdName, exitType)

	n := 100
	var reports []database.Report
	q := getReportsQuery(hostname, cmdName, exitType, exitCode)
	q = q.Limit(n)
	q.Find(&reports)

	dateTimes := make([]string, len(reports))
	codes := make([]int64, len(reports))
	uids := make([]string, len(reports))
	failures := 0
	for i, r := range reports {
		dateTimes[len(reports)-i-1] = r.EndTime.UTC().Format("2006-01-02 15:04:05.000 MST")
		codes[len(reports)-i-1] = int64(r.ExitCode)
		uids[len(reports)-i-1] = r.Ulid
		if r.ExitCode != 0 {
			failures++
		}
	}

	failPercent := 0.0
	if len(reports) > 0 {
		failPercent = 100 * (1.0 - (float64(failures) / float64(len(reports))))
	}

	ctx.MustRender("reports/graph.html", struct {
		Title          string
		ReportListURL  string
		FailPercent    float64
		DateTimes      []string
		Codes          []int64
		Ulids          []string
		HasFilters     bool
		FilterHostname string
		FilterCmdName  string
		FilterExitType string
	}{
		"Graph",
		"/reports" + paginationReadyQueryString(activeQuery),
		failPercent,
		dateTimes,
		codes,
		uids,
		(hostname != "" || cmdName != "" || exitType != ""),
		hostname,
		cmdName,
		exitType,
	})
}