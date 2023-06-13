package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

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
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	err = dbs.db.PingContext(ctxWithTimeout)
	return
}

func (dbs dbStorage) InitializeMetrics(ctx context.Context, restore *bool) (err error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v(id SERIAL PRIMARY KEY, "+
		"metric_name VARCHAR(100), metric_type VARCHAR(7), metric_value DOUBLE PRECISION, metric_delta INT)", MetricsTableName)
	_, err = dbs.db.ExecContext(ctxWithTimeout, query)
	if err != nil {
		return
	}

	if !*restore {
		return
	}

	//load metrics from db to memstorage
	query = fmt.Sprintf("SELECT metric_name, metric_type, metric_value, metric_delta FROM %v", MetricsTableName)
	rows, err := dbs.db.QueryContext(ctxWithTimeout, query)
	if err != nil {
		return
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

	dbs.mutex.Lock()

	tx, err := dbs.db.Begin()
	if err != nil {
		return
	}

	//delete all records and insert anew
	//brutal but more effective than checking for each metric if it exists and update
	//or insert individually

	_, err = tx.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %v", MetricsTableName))
	if err != nil {
		_ = tx.Rollback()
		return
	}
	for _, value := range *metricsMap {
		var insertValues string
		switch value.MType {
		case metric.GaugeMetricType:
			insertValues = fmt.Sprintf("DEFAULT, '%v', '%v', %v, %v", value.ID, value.MType, *value.Value, "NULL")
		case metric.CounterMetricType:
			insertValues = fmt.Sprintf("DEFAULT, '%v', '%v', %v, %v", value.ID, value.MType, "NULL", *value.Delta)
		}

		if insertValues != "" {
			query := fmt.Sprintf("INSERT INTO %v (id, metric_name, metric_type, metric_value, metric_delta) VALUES (%v)", MetricsTableName, insertValues)
			fmt.Println(query)
			_, err = tx.ExecContext(ctx, query)
			if err != nil {
				_ = tx.Rollback()
				return
			}
		}
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
