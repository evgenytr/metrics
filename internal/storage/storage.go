package storage

import (
	"database/sql"

	"github.com/evgenytr/metrics.git/internal/storage/database"
	"github.com/evgenytr/metrics.git/internal/storage/memstorage"
)

type Storage interface {
	LoadMetrics() error
	StoreMetrics() error
	UpdateGauge(name string, value *float64) (*float64, error)
	UpdateCounter(name string, value *int64) (*int64, error)
	GetGaugeValue(name string) (*float64, error)
	GetCounterValue(name string) (*int64, error)
	Update(metricType, name, value string) (string, error)
	ReadValue(metricType, name string) (string, error)
	ListAll() (map[string]string, error)
	Ping() error
}

func NewStorage(db *sql.DB, fileStoragePath *string) Storage {
	if db != nil {
		return database.NewStorage(db)
	}
	return memstorage.NewStorage(fileStoragePath)
}
