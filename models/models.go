package models

import (
	"time"
)

type Tag struct {
	ID   uint   `gorm:"primary_key"`
	Text string `gorm:"text;type:varchar(255);unique_index"`
}

type Report struct {
	ID              uint       `gorm:"primary_key"  json:"-"`
	Ulid            string     `gorm:"type:varchar(100);unique_index" json:"-"`
	CreatedAt       time.Time  `json:"-"`
	UpdatedAt       time.Time  `json:"-"`
	DeletedAt       *time.Time `sql:"index" json:"-"`
	Name            string     `json:"name"`
	Command         string     `json:"-"`
	RawCommand      []string   `gorm:"-" json:"command"`
	StartTime       time.Time  `json:"start_time"`
	EndTime         time.Time  `json:"end_time"`
	ElapsedSeconds  float32    `json:"elapsed_seconds"`
	ExitCode        int        `json:"exit_code"`
	ExitDescription string     `json:"exit_description"`
	CapturedOutput  string     `json:"captured_output"`
	Hostname        string     `json:"hostname"`
	Tags            []Tag      `json:"-" gorm:"many2many:report_tags;"`
	RawTags         []string   `gorm:"-" json:"tags"`
}
