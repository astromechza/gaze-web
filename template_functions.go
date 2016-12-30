package main

import (
	"time"

	"github.com/ararog/timeago"
)

func templateFuncAdd(i int64, j int64) int64 {
	return i + j
}

func templateFuncIRange(count int64) []int64 {
	output := make([]int64, count)
	for i, _ := range output {
		output[i] = int64(i)
	}
	return output
}

func timestampFormat1(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.000 MST")
}

func timeAgoString(t time.Time) string {
	s, err := timeago.TimeAgoFromNowWithTime(t)
	if err != nil {
		s = "(unknown)"
	}
	return s
}

func buildTemplateFuncsMap() map[string]interface{} {
	output := make(map[string]interface{})
	output["add"] = templateFuncAdd
	output["irange"] = templateFuncIRange
	output["tsf"] = timestampFormat1
	output["tago"] = timeAgoString
	return output
}
