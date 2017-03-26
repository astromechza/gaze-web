package models

import (
	"time"
)

type Report struct {
	Ulid            string    `json:"ulid"`
	Name            string    `json:"name"`
	Command         []string  `json:"command"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	ElapsedSeconds  float32   `json:"elapsed_seconds"`
	ExitCode        int       `json:"exit_code"`
	ExitDescription string    `json:"exit_description"`
	CapturedOutput  string    `json:"captured_output"`
	Hostname        string    `json:"hostname"`
	Tags            []string  `json:"tags"`
}
