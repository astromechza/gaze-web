package gorm

import (
	"fmt"
	"strings"

	"github.com/AstromechZA/gaze-web/models"
	"github.com/AstromechZA/gaze-web/storage"

	"github.com/jinzhu/gorm"

	// support for sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type GormReportStore struct {
	active *gorm.DB
}

func SetupGormReportStore(filepath string) (*GormReportStore, error) {
	store := new(GormReportStore)

	db, err := gorm.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Report{})
	db.AutoMigrate(&models.Tag{})

	store.active = db

	return store, nil
}

func (s *GormReportStore) AddReport(report *models.Report) error {
	if len(report.RawTags) > 0 {
		report.Tags = make([]models.Tag, 0)
		for _, ts := range report.RawTags {
			ts = strings.TrimSpace(strings.ToLower(ts))
			if len(ts) > 0 {
				// validate
				var tag models.Tag
				if err := s.active.Where("text = ?", ts).First(&tag).Error; err != nil {
					tag.Text = ts
					err = s.active.Create(&tag).Error
					if err != nil {
						return fmt.Errorf("db create tag error %s", err)
					}
				}
				report.Tags = append(report.Tags, tag)
			}
		}
	}
	return s.active.Create(&report).Error
}

func (s *GormReportStore) GetReport(u string) (*models.Report, error) {
	report := new(models.Report)
	err := s.active.Preload("Tags").Where("ulid = ?", u).First(report).Error
	return report, err
}

func (s *GormReportStore) getReportsQuery(filter storage.ReportStoreFilter) *gorm.DB {
	q := s.active.Table("reports").Order("start_time desc")
	if filter.Hostname != "" {
		q = q.Where("hostname = ?", filter.Hostname)
	}
	if filter.Cmdname != "" {
		q = q.Where("name = ?", filter.Cmdname)
	}
	if filter.ExitType != "" {
		q = q.Where("exit_code = ?", filter.ExitCode)
	}
	return q
}

func (s *GormReportStore) CountReports(filter storage.ReportStoreFilter) (int, error) {
	var totalRecords int64
	err := s.getReportsQuery(filter).Count(&totalRecords).Error
	return int(totalRecords), err
}

func (s *GormReportStore) GetLatestFailedReport(filter storage.ReportStoreFilter) (*models.Report, error) {
	lastFailure := new(models.Report)
	q := s.getReportsQuery(filter)
	q = q.Where("exit_code != ?", 0)
	q = q.First(&lastFailure)
	return lastFailure, q.Error
}

func (s *GormReportStore) GetLatestSuccessfulReport(filter storage.ReportStoreFilter) (*models.Report, error) {
	lastFailure := new(models.Report)
	q := s.getReportsQuery(filter)
	q = q.Where("exit_code = ?", 0)
	q = q.First(&lastFailure)
	return lastFailure, q.Error
}

func (s *GormReportStore) ListReportsPage(filter storage.ReportStoreFilter, numberPerPage, pageNum int) (*[]models.Report, error) {
	var reports []models.Report
	q := s.getReportsQuery(filter)
	q = q.Offset((pageNum - 1) * numberPerPage).Limit(numberPerPage)
	q = q.Find(&reports)
	return &reports, q.Error
}

func (s *GormReportStore) Close() error {
	return s.active.Close()
}

var _ storage.GazeWebReportStore = (*GormReportStore)(nil)
