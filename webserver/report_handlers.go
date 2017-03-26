package webserver

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/url"
	"strings"
	"time"

	"strconv"

	"github.com/AstromechZA/gaze-web/models"
	"github.com/AstromechZA/gaze-web/random"
	"github.com/AstromechZA/gaze-web/storage"

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
	var report models.Report
	err := ctx.ReadJSON(&report)
	if err != nil {
		failJSON(ctx, 400, "Invalid Report payload "+err.Error())
		return
	}

	nowTime := time.Now()

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
	if report.ElapsedSeconds < 0 {
		failJSON(ctx, 400, "'elapsed_seconds' value must be >= 0")
		return
	}

	// fields that are optional
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

	if err = storage.ActiveStore.AddReport(&report); err != nil {
		ctx.Log(iris.DevMode, "db create report error "+err.Error())
		failJSON(ctx, 500, "Server Error")
	} else {
		ctx.SetStatusCode(204)
	}
}

func parseFilterParams(ctx *iris.Context) (string, string, string, int) {
	hostname := strings.TrimSpace(ctx.URLParam("hostname"))
	cmdName := strings.TrimSpace(ctx.URLParam("name"))
	exitType := strings.TrimSpace(ctx.URLParam("exit"))
	exitType = strings.ToLower(exitType)
	exitCode, err := strconv.ParseInt(exitType, 10, 32)
	if err != nil {
		if exitType == "zero" || exitType == "success" {
			exitCode = 0
		} else if exitType == "nonzero" || exitType == "failure" {
			exitCode = -1
		} else if exitType == "any" || exitType == "anything" {
			exitType = ""
		}
	}
	return hostname, cmdName, exitType, int(exitCode)
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

func parsePaginationParams(ctx *iris.Context) (int, int) {
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
	return int(pageNum), int(numberPerPage)
}

func ListReportsHandler(ctx *iris.Context) {
	pageNum, numberPerPage := parsePaginationParams(ctx)

	hostname, cmdName, exitType, exitCode := parseFilterParams(ctx)
	activeQuery := buildURLValues(hostname, cmdName, exitType)

	filter := storage.ReportStoreFilter{
		Hostname: hostname,
		Name:     cmdName,
		ExitCode: exitCode,
		ExitType: exitType,
	}

	reports, err := storage.ActiveStore.ListReportsPage(filter, numberPerPage, pageNum)
	if err != nil {
		ctx.Log(iris.DevMode, "err: %s", err)
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}

	totalRecords, err := storage.ActiveStore.CountReports(filter)
	if err != nil {
		ctx.Log(iris.DevMode, "err: %s", err)
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}

	lastFailure, err := storage.ActiveStore.GetLatestFailedReport(filter)
	if err != nil {
		ctx.Log(iris.DevMode, "err: %s", err)
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}
	lastSuccess, err := storage.ActiveStore.GetLatestSuccessfulReport(filter)
	if err != nil {
		ctx.Log(iris.DevMode, "err: %s", err)
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}

	numberOfPages := int(math.Ceil(float64(totalRecords) / float64(numberPerPage)))

	ctx.MustRender("reports/list.html", struct {
		Title        string
		Reports      []models.Report
		LastSuccess  *models.Report
		LastFailure  *models.Report
		TotalRecords int
		GraphLink    string

		CurrentPage   int64
		TotalPages    int64
		URLToPaginate string

		FormName     string
		FormHostname string
		FormExit     string
	}{
		"Reports",
		*reports, lastSuccess, lastFailure,
		totalRecords,
		"/graph" + paginationReadyQueryString(activeQuery),
		int64(pageNum), int64(numberOfPages),
		"/reports" + paginationReadyQueryString(activeQuery),
		cmdName,
		hostname,
		exitType,
	})
}

func GetReportHandler(ctx *iris.Context) {
	u := ctx.Param("ulid")
	report, err := storage.ActiveStore.GetReport(u)
	if err != nil {
		ctx.EmitError(iris.StatusNotFound)
		return
	}
	ctx.MustRender("reports/show.html", struct {
		Title  string
		Report models.Report
	}{report.Ulid, *report})
}

func GraphReportsHandler(ctx *iris.Context) {
	hostname, cmdName, exitType, exitCode := parseFilterParams(ctx)
	activeQuery := buildURLValues(hostname, cmdName, exitType)

	filter := storage.ReportStoreFilter{
		Hostname: hostname,
		Name:     cmdName,
		ExitCode: exitCode,
		ExitType: exitType,
	}

	reports, _ := storage.ActiveStore.ListReportsPage(filter, 100, 1)
	reports2 := *reports

	dateTimes := make([]string, len(reports2))
	codes := make([]int64, len(reports2))
	uids := make([]string, len(reports2))
	failures := 0
	for i, r := range reports2 {
		dateTimes[len(reports2)-i-1] = r.EndTime.UTC().Format("2006-01-02 15:04:05.000 MST")
		codes[len(reports2)-i-1] = int64(r.ExitCode)
		uids[len(reports2)-i-1] = r.Ulid
		if r.ExitCode != 0 {
			failures++
		}
	}

	failPercent := 0.0
	if len(reports2) > 0 {
		failPercent = 100 * (1.0 - (float64(failures) / float64(len(reports2))))
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
