// Package memstorage is used to store metrics info in memory
// and periodically save it to file
package memstorage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/evgenytr/metrics.git/internal/interfaces"
	"github.com/evgenytr/metrics.git/internal/metric"
)

type memStorage struct {
	metricsMap      map[string]*metric.Metrics
	fileStoragePath string
}

// NewStorage returns pointer to memStorage struct.
func NewStorage(fileStoragePath string) interfaces.Storage {
	return &memStorage{
		metricsMap:      make(map[string]*metric.Metrics),
		fileStoragePath: fileStoragePath,
	}
}

// Ping makes no sense for memstorage, this is a stub to comply to interface.
func (ms memStorage) Ping(_ context.Context) (err error) {
	fmt.Println("ping memstorage")
	err = fmt.Errorf("no database used")
	return
}

// InitializeMetrics loads saved to file metric values.
func (ms memStorage) InitializeMetrics(_ context.Context, restore bool) (err error) {
	fmt.Println("init metrics")
	if !restore {
		return
	}
	if ms.fileStoragePath == "" {
		err = fmt.Errorf("no file storage path")
		return
	}
	data, err := os.ReadFile(ms.fileStoragePath)
	if err != nil {
		return
	}

	var metricsMap = make(map[string]*metric.Metrics)

	if err = json.Unmarshal(data, &metricsMap); err != nil {
		return
	}
	//TODO: validate metrics loaded from file
	for key, value := range metricsMap {
		ms.metricsMap[key] = value
	}

	return
}

// StoreMetrics saves metrics values to file.
func (ms memStorage) StoreMetrics(_ context.Context) (err error) {
	fmt.Println("store metrics")
	if ms.fileStoragePath == "" {
		err = fmt.Errorf("no file storage path")
		return
	}
	jsonRes, err := json.Marshal(ms.metricsMap)
	if err != nil {
		return
	}
	err = os.WriteFile(ms.fileStoragePath, jsonRes, 0666)
	return
}

// Update updates value of metric in memory.
func (ms memStorage) Update(_ context.Context, metricType, name, value string) (newValue string, err error) {
	if currMetric, ok := ms.metricsMap[name]; ok {
		newValue, err = currMetric.Add(metricType, value)
		if err != nil {
			return
		}
	} else {
		ms.metricsMap[name], err = metric.Create(metricType, name, value)
		newValue = value
	}
	return
}

// UpdateGauge updates gauge metric value.
func (ms memStorage) UpdateGauge(_ context.Context, name string, value *float64) (newValue *float64, err error) {
	if currMetric, ok := ms.metricsMap[name]; ok {
		newValue, err = currMetric.UpdateGauge(value)
		if err != nil {
			return
		}
	} else {
		ms.metricsMap[name], err = metric.CreateGauge(name, value)
		newValue = value
	}
	return
}

// UpdateCounter updates counter metric value.
func (ms memStorage) UpdateCounter(_ context.Context, name string, value *int64) (newValue *int64, err error) {
	if currMetric, ok := ms.metricsMap[name]; ok {
		newValue, err = currMetric.UpdateCounter(value)
		if err != nil {
			return
		}
	} else {
		ms.metricsMap[name], err = metric.CreateCounter(name, value)
		newValue = value
	}
	return
}

func (ms memStorage) ReadValue(_ context.Context, metricType, name string) (value string, err error) {
	if currMetric, ok := ms.metricsMap[name]; ok {
		if currMetric.GetType() != metricType {
			err = fmt.Errorf("metric type mismatch")
			return
		}
		value = currMetric.GetValue()
	} else {
		err = fmt.Errorf("metric not found")
	}
	return
}

func (ms memStorage) GetGaugeValue(_ context.Context, name string) (value *float64, err error) {
	if currMetric, ok := ms.metricsMap[name]; ok {
		if currMetric.GetType() != metric.GaugeMetricType {
			err = fmt.Errorf("metric type mismatch")
			return
		}
		value = currMetric.GetGaugeValue()
	} else {
		err = fmt.Errorf("metric not found")
	}
	return
}

func (ms memStorage) GetCounterValue(_ context.Context, name string) (value *int64, err error) {
	if currMetric, ok := ms.metricsMap[name]; ok {
		if currMetric.GetType() != metric.CounterMetricType {
			err = fmt.Errorf("metric type mismatch")
			return
		}
		value = currMetric.GetCounterValue()
	} else {
		err = fmt.Errorf("metric not found")
	}
	return
}

func (ms memStorage) ListAll(_ context.Context) (metricsMap *map[string]string, err error) {
	newMetricsMap := make(map[string]string, len(ms.metricsMap))
	for key, value := range ms.metricsMap {
		newMetricsMap[key] = value.GetValue()
	}
	return &newMetricsMap, err
}

func (ms memStorage) GetMetricsMap(_ context.Context) (metricsMap *map[string]*metric.Metrics, err error) {
	newMetricsMap := make(map[string]*metric.Metrics, len(ms.metricsMap))
	for key, value := range ms.metricsMap {
		newMetricsMap[key] = value
	}
	return &newMetricsMap, err
}
