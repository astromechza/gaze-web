package bolt

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/boltdb/bolt"

	"github.com/AstromechZA/gaze-web/models"
	"github.com/AstromechZA/gaze-web/storage"
)

const (
	bucketReports          = "Model.Reports"
	bucketIndexReportsUlid = "Indexes.ReportByUlid"
)

type BoltDBReportStore struct {
	db *bolt.DB
}

func SetupBoltDBReportStore(filepath string) (*BoltDBReportStore, error) {
	store := new(BoltDBReportStore)

	db, err := bolt.Open(filepath, 0644, nil)
	if err != nil {
		return nil, err
	}
	store.db = db

	err = store.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(bucketReports)); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(bucketIndexReportsUlid)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create initial buckets")
	}
	return store, nil
}

func (s *BoltDBReportStore) AddReport(report *models.Report) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		mainBucket := tx.Bucket([]byte(bucketReports))
		payload, err := json.Marshal(*report)
		if err != nil {
			return err
		}
		return mainBucket.Put([]byte(report.Ulid), payload)
	})
	return err
}

func (s *BoltDBReportStore) GetReport(u string) (*models.Report, error) {
	var output models.Report
	err := s.db.View(func(tx *bolt.Tx) error {
		mainBucket := tx.Bucket([]byte(bucketReports))
		payload := mainBucket.Get([]byte(u))
		if payload == nil || len(payload) == 0 {
			return fmt.Errorf("does not exist")
		}
		return json.Unmarshal(payload, &output)
	})
	if err != nil {
		return nil, err
	}
	return &output, nil
}

func (s *BoltDBReportStore) getAllReports(filter storage.ReportStoreFilter) (*[]models.Report, error) {
	var output []models.Report
	err := s.db.View(func(tx *bolt.Tx) error {
		mainBucket := tx.Bucket([]byte(bucketReports))
		c := mainBucket.Cursor()
		for k, payload := c.First(); k != nil; k, payload = c.Next() {
			var temp models.Report
			_ = json.Unmarshal(payload, &temp)

			if filter.Name != "" && temp.Name != filter.Name {
				continue
			}

			if filter.Hostname != "" && temp.Hostname != filter.Hostname {
				continue
			}

			if filter.ExitType == "not" && temp.ExitCode == filter.ExitCode {
				continue
			} else if filter.ExitType != "" && temp.ExitCode != filter.ExitCode {
				continue
			}

			output = append(output, temp)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &output, nil
}

type reportSortRevByStartTime []models.Report

func (s reportSortRevByStartTime) Len() int {
	return len(s)
}
func (s reportSortRevByStartTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s reportSortRevByStartTime) Less(i, j int) bool {
	return s[i].StartTime.After(s[j].StartTime)
}

func (s *BoltDBReportStore) CountReports(filter storage.ReportStoreFilter) (int, error) {
	var output int
	err := s.db.View(func(tx *bolt.Tx) error {
		mainBucket := tx.Bucket([]byte(bucketReports))
		c := mainBucket.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			output++
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return output, nil
}

func (s *BoltDBReportStore) GetLatestFailedReport(filter storage.ReportStoreFilter) (*models.Report, error) {
	// TODO can optimise by making explicit query
	filter.ExitType = "not"
	filter.ExitCode = 0
	reports, err := s.getAllReports(filter)
	if err != nil {
		return nil, err
	}
	reports2 := *reports
	if len(reports2) == 0 {
		return nil, nil
	}
	sort.Sort(reportSortRevByStartTime(reports2))
	lastReport := reports2[0]
	return &lastReport, nil
}

func (s *BoltDBReportStore) GetLatestSuccessfulReport(filter storage.ReportStoreFilter) (*models.Report, error) {
	// TODO can optimise by making explicit query
	filter.ExitType = "equal"
	filter.ExitCode = 0
	reports, err := s.getAllReports(filter)
	if err != nil {
		return nil, err
	}
	reports2 := *reports
	if len(reports2) == 0 {
		return nil, nil
	}
	sort.Sort(reportSortRevByStartTime(reports2))
	lastReport := reports2[0]
	return &lastReport, nil
}

func (s *BoltDBReportStore) ListReportsPage(filter storage.ReportStoreFilter, numberPerPage, pageNum int) (*[]models.Report, error) {
	var empty []models.Report
	if numberPerPage == 0 {
		return &empty, nil
	}
	reports, err := s.getAllReports(filter)
	if err != nil {
		return nil, err
	}
	reports2 := *reports
	if len(reports2) == 0 {
		return &empty, nil
	}
	sort.Sort(reportSortRevByStartTime(reports2))
	startIndex := (pageNum - 1) * numberPerPage
	if startIndex < 0 {
		startIndex = 0
	} else if startIndex >= len(reports2) {
		return &empty, nil
	}
	endIndex := startIndex + numberPerPage
	if endIndex >= len(reports2) {
		endIndex = len(reports2) - 1
	}
	output := make([]models.Report, endIndex-startIndex+1)
	slice := reports2[startIndex:(endIndex + 1)]
	copy(output, slice)
	return &output, nil
}

func (s *BoltDBReportStore) Close() error {
	return s.db.Close()
}

var _ storage.GazeWebReportStore = (*BoltDBReportStore)(nil)
