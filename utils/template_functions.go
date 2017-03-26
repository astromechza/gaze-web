package utils

import (
	"time"

	"strings"

	"github.com/ararog/timeago"
)

func templateFuncAdd(i int64, j int64) int64 {
	return i + j
}

func templateFuncIRange(count int64) []int64 {
	output := make([]int64, count)
	for i := range output {
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

func quoteArgItem(item string) string {
	if strings.ContainsAny(item, " \t\n\v\"") {
		buff := "\""

		for i := 0; ; i++ {
			numescapes := 0
			for ; i < len(item) && item[i] == '\\'; i++ {
				numescapes++
			}
			if i >= len(item) {
				for x := 0; x < numescapes; x++ {
					buff += "\\\\"
				}
				break
			} else if item[i] == '"' {
				for x := 0; x < numescapes; x++ {
					buff += "\\\\"
				}
				buff += "\\"
				buff += string(item[i])
			} else {
				for x := 0; x < numescapes; x++ {
					buff += "\\"
				}
				buff += string(item[i])
			}
		}
		buff += "\""
		return buff
	}
	return item
}

func formatListAsQuoted(argv []string) string {
	buff := ""
	for i, item := range argv {
		if i > 0 {
			buff += " "
		}
		buff += quoteArgItem(item)
	}
	return buff
}

var EmbeddedVersionString = ""

func BuildTemplateFuncsMap() map[string]interface{} {
	output := make(map[string]interface{})
	output["add"] = templateFuncAdd
	output["irange"] = templateFuncIRange
	output["tsf"] = timestampFormat1
	output["fet"] = formatElapsedTime
	output["tago"] = timeAgoString
	output["version"] = func() string { return EmbeddedVersionString }
	output["formatListAsQuoted"] = formatListAsQuoted
	return output
}
