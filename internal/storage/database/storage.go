package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/evgenytr/metrics.git/internal/interfaces"
	"github.com/evgenytr/metrics.git/internal/metric"
	"github.com/evgenytr/metrics.git/internal/storage/memstorage"
)

const MetricsTableName = "metrics"

type dbStorage struct {
	db    *sql.DB
	ms    interfaces.Storage
	mutex *sync.Mutex
}

func NewStorage(db *sql.DB) interfaces.Storage {
	return &dbStorage{
		db:    db,
		ms:    memstorage.NewStorage(nil),
		mutex: &sync.Mutex{},
	}
}

func (dbs dbStorage) Ping(ctx context.Context) (err error) {
	fmt.Println("ping dbstorage")

	const pingTimeout = 1
	ctxWithTimeout, cancel := context.WithTimeout(ctx, pingTimeout*time.Second)
	defer cancel()

	err = dbs.db.PingContext(ctxWithTimeout)
	return
}

func (dbs dbStorage) InitializeMetrics(ctx context.Context, restore *bool) (err error) {

	const databaseExecTimeout = 3
	ctxWithTimeout, cancel := context.WithTimeout(ctx, databaseExecTimeout*time.Second)
	defer cancel()

	//run migrations

	const migrationsPath = "./internal/storage/database/migration"

	absPath, err := filepath.Abs(migrationsPath)

	fmt.Println(absPath)

	//TODO: use project database instead of 'postgres', figure out how and when to create it, get name from config?
	driver, err := postgres.WithInstance(dbs.db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", absPath), "postgres", driver)
	if err != nil {
		return fmt.Errorf("Failed to get a new migrate instance: %w", err)
	}

	version, dirty, err := m.Version()
	fmt.Println(version, dirty, err)

	if err := m.Up(); err != nil {
		fmt.Println(err)
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("Failed to apply migrations to DB: %w", err)
		}
	}

	if !*restore {
		return
	}

	//load metrics from db to memstorage
	query := fmt.Sprintf("SELECT metric_name, metric_type, metric_value, metric_delta FROM %v", MetricsTableName)
	rows, err := dbs.db.QueryContext(ctxWithTimeout, query)
	if err != nil {
		return fmt.Errorf("Failed to load metrics from db: %w", err)
	}
	err = rows.Err()
	if err != nil {
		return
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	for rows.Next() {
		var metricName string
		var metricType string
		var metricValue float64
		var metricDelta int64

		if err = rows.Scan(&metricName, &metricType, &metricValue, &metricDelta); err != nil {
			return
		}

		switch metricType {
		case metric.GaugeMetricType:
			_, err = dbs.ms.UpdateGauge(ctx, metricName, &metricValue)
		case metric.CounterMetricType:
			_, err = dbs.ms.UpdateCounter(ctx, metricName, &metricDelta)
		default:
			err = fmt.Errorf("metric type not supported")
		}
		if err != nil {
			return
		}

	}
	return
}

func (dbs dbStorage) StoreMetrics(ctx context.Context) (err error) {
	//store metrics from memstorage to db
	fmt.Println("store metrics to db")
	metricsMap, err := dbs.ms.GetMetricsMap(ctx)
	if err != nil {
		return
	}

	//nothing to store
	if len(*metricsMap) == 0 {
		return
	}

	insertValues := make([]string, len(*metricsMap))
	valuesIndex := 0

	for _, value := range *metricsMap {
		switch value.MType {
		case metric.GaugeMetricType:
			currInsertValues := fmt.Sprintf("(DEFAULT, '%v', '%v', %v, %v)", value.ID, value.MType, *value.Value, "NULL")
			insertValues[valuesIndex] = currInsertValues
		case metric.CounterMetricType:
			currInsertValues := fmt.Sprintf("(DEFAULT, '%v', '%v', %v, %v)", value.ID, value.MType, "NULL", *value.Delta)
			insertValues[valuesIndex] = currInsertValues
		}
		valuesIndex++
	}

	dbs.mutex.Lock()

	tx, err := dbs.db.Begin()
	if err != nil {
		return
	}

	query := fmt.Sprintf("INSERT INTO %s (id, metric_name, metric_type, metric_value, metric_delta) VALUES %v "+
		"ON CONFLICT (metric_name) DO UPDATE SET metric_value = EXCLUDED.metric_value, metric_delta = EXCLUDED.metric_delta",
		MetricsTableName, strings.Join(insertValues, ", "))
	fmt.Println(query)
	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		_ = tx.Rollback()
		return
	}

	err = tx.Commit()

	dbs.mutex.Unlock()

	return
}

func (dbs dbStorage) Update(ctx context.Context, metricType, name, value string) (newValue string, err error) {
	newValue, err = dbs.ms.Update(ctx, metricType, name, value)
	return
}

func (dbs dbStorage) UpdateGauge(ctx context.Context, name string, value *float64) (newValue *float64, err error) {
	newValue, err = dbs.ms.UpdateGauge(ctx, name, value)
	return
}

func (dbs dbStorage) UpdateCounter(ctx context.Context, name string, value *int64) (newValue *int64, err error) {
	newValue, err = dbs.ms.UpdateCounter(ctx, name, value)
	return
}

func (dbs dbStorage) ReadValue(ctx context.Context, metricType, name string) (value string, err error) {
	value, err = dbs.ms.ReadValue(ctx, metricType, name)
	return
}

func (dbs dbStorage) GetGaugeValue(ctx context.Context, name string) (value *float64, err error) {
	value, err = dbs.ms.GetGaugeValue(ctx, name)
	return
}

func (dbs dbStorage) GetCounterValue(ctx context.Context, name string) (value *int64, err error) {
	value, err = dbs.ms.GetCounterValue(ctx, name)
	return
}

func (dbs dbStorage) ListAll(ctx context.Context) (metricsMap *map[string]string, err error) {
	metricsMap, err = dbs.ms.ListAll(ctx)
	return
}

func (dbs dbStorage) GetMetricsMap(ctx context.Context) (metricsMap *map[string]*metric.Metrics, err error) {
	metricsMap, err = dbs.ms.GetMetricsMap(ctx)
	return
}
