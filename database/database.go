package database

import (
	"time"

	"github.com/jinzhu/gorm"

	// support for sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var ActiveDB *gorm.DB

func InitSqliteDatabase(filepath string) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Report{})
	db.AutoMigrate(&Tag{})

	return db, nil
}

type Tag struct {
	ID   uint   `gorm:"primary_key"`
	Text string `gorm:"text;type:varchar(255);unique_index"`
}

type Report struct {
	ID              uint       `gorm:"primary_key"  json:"-"`
	CreatedAt       time.Time  `json:"-"`
	UpdatedAt       time.Time  `json:"-"`
	DeletedAt       *time.Time `sql:"index" json:"-"`
	Ulid            string     `gorm:"type:varchar(100);unique_index" json:"-"`
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
