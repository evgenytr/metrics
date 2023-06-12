package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
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

type dbStorage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) Storage {
	return &dbStorage{
		db: db,
	}
}

func (dbs dbStorage) Ping() (err error) {
	fmt.Println("ping dbstorage")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = dbs.db.PingContext(ctx)
	return
}

func (dbs dbStorage) LoadMetrics() (err error) {

	return
}

func (dbs dbStorage) StoreMetrics() (err error) {

	return
}

func (dbs dbStorage) Update(metricType, name, value string) (newValue string, err error) {

	return
}

func (dbs dbStorage) UpdateGauge(name string, value *float64) (newValue *float64, err error) {

	return
}

func (dbs dbStorage) UpdateCounter(name string, value *int64) (newValue *int64, err error) {

	return
}

func (dbs dbStorage) ReadValue(metricType, name string) (value string, err error) {

	return
}

func (dbs dbStorage) GetGaugeValue(name string) (value *float64, err error) {

	return
}

func (dbs dbStorage) GetCounterValue(name string) (value *int64, err error) {

	return
}

func (dbs dbStorage) ListAll() (metricsMap map[string]string, err error) {

	return
}
