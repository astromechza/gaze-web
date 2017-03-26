package noop

import (
	"fmt"

	"github.com/AstromechZA/gaze-web/models"
	"github.com/AstromechZA/gaze-web/storage"
)

type NoopReportStore struct {
}

func (s *NoopReportStore) AddReport(report *models.Report) error {
	return fmt.Errorf("Not implemented")
}

func (s *NoopReportStore) GetReport(u string) (*models.Report, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *NoopReportStore) CountReports(filter storage.ReportStoreFilter) (int, error) {
	return 0, fmt.Errorf("Not implemented")
}

func (s *NoopReportStore) GetLatestFailedReport(filter storage.ReportStoreFilter) (*models.Report, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *NoopReportStore) GetLatestSuccessfulReport(filter storage.ReportStoreFilter) (*models.Report, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *NoopReportStore) ListReportsPage(filter storage.ReportStoreFilter, numberPerPage, pageNum int) (*[]models.Report, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *NoopReportStore) Close() error {
	return nil
}

var _ storage.GazeWebReportStore = (*NoopReportStore)(nil)
