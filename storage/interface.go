package storage

import "github.com/AstromechZA/gaze-web/models"

var ActiveStore GazeWebReportStore

type ReportStoreFilter struct {
	Hostname string
	Name     string
	ExitCode int
	ExitType string
}

// GazeWebReportStore describes a storage interface for reports and for querying for
// different types of related data
type GazeWebReportStore interface {
	// AddReport adds a new report to the database
	AddReport(report *models.Report) error

	// Lookup a report by its identifier
	GetReport(ulid string) (*models.Report, error)

	// Count the reports that match the filter
	CountReports(filter ReportStoreFilter) (int, error)

	// Get the latest report that failed
	GetLatestFailedReport(filter ReportStoreFilter) (*models.Report, error)

	// Get the latest report that succeeded
	GetLatestSuccessfulReport(filter ReportStoreFilter) (*models.Report, error)

	// Lookup a paginated section of reports that match the filter
	ListReportsPage(filter ReportStoreFilter, numberPerPage, pageNum int) (*[]models.Report, error)

	// Close the store and clean up if necessary
	Close() error
}
