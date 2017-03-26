package bolt

import (
	"fmt"

	"github.com/boltdb/bolt"

	"github.com/AstromechZA/gaze-web/models"
	"github.com/AstromechZA/gaze-web/storage"
)

const (
	bucketReports = "Model.Reports"
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
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create initial buckets")
	}
	return store, nil
}

func (s *BoltDBReportStore) AddReport(report *models.Report) error {
	return fmt.Errorf("Not implemented")
}

func (s *BoltDBReportStore) GetReport(u string) (*models.Report, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *BoltDBReportStore) CountReports(filter storage.ReportStoreFilter) (int, error) {
	return 0, fmt.Errorf("Not implemented")
}

func (s *BoltDBReportStore) GetLatestFailedReport(filter storage.ReportStoreFilter) (*models.Report, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *BoltDBReportStore) GetLatestSuccessfulReport(filter storage.ReportStoreFilter) (*models.Report, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *BoltDBReportStore) ListReportsPage(filter storage.ReportStoreFilter, numberPerPage, pageNum int) (*[]models.Report, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *BoltDBReportStore) Close() error {
	return s.database.Close()
}

var _ storage.GazeWebReportStore = (*BoltDBReportStore)(nil)
