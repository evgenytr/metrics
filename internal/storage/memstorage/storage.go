package memstorage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/evgenytr/metrics.git/internal/metric"
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

type memStorage struct {
	metricsMap      map[string]*metric.Metrics
	fileStoragePath *string
}

func NewStorage(fileStoragePath *string) Storage {
	return &memStorage{
		metricsMap:      make(map[string]*metric.Metrics),
		fileStoragePath: fileStoragePath,
	}
}

func (ms memStorage) Ping() (err error) {
	fmt.Println("ping memstorage")
	err = fmt.Errorf("no database used")
	return
}

func (ms memStorage) LoadMetrics() (err error) {
	fmt.Println("load metrics")
	if ms.fileStoragePath == nil || *ms.fileStoragePath == "" {
		err = fmt.Errorf("no file storage path")
		return
	}
	data, err := os.ReadFile(*ms.fileStoragePath)
	if err != nil {
		return
	}

	var metricsMap = make(map[string]*metric.Metrics)

	if err = json.Unmarshal(data, &metricsMap); err != nil {
		return
	}
	//TODO: validate
	for key, value := range metricsMap {
		ms.metricsMap[key] = value
	}

	return
}

func (ms memStorage) StoreMetrics() (err error) {
	fmt.Println("store metrics")
	if ms.fileStoragePath == nil || *ms.fileStoragePath == "" {
		err = fmt.Errorf("no file storage path")
		return
	}
	jsonRes, err := json.Marshal(ms.metricsMap)
	if err != nil {
		return
	}
	err = os.WriteFile(*ms.fileStoragePath, jsonRes, 0666)
	return
}

func (ms memStorage) Update(metricType, name, value string) (newValue string, err error) {
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

func (ms memStorage) UpdateGauge(name string, value *float64) (newValue *float64, err error) {
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

func (ms memStorage) UpdateCounter(name string, value *int64) (newValue *int64, err error) {
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

func (ms memStorage) ReadValue(metricType, name string) (value string, err error) {
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

func (ms memStorage) GetGaugeValue(name string) (value *float64, err error) {
	if currMetric, ok := ms.metricsMap[name]; ok {
		if currMetric.GetType() != "gauge" {
			err = fmt.Errorf("metric type mismatch")
			return
		}
		value = currMetric.GetGaugeValue()
	} else {
		err = fmt.Errorf("metric not found")
	}
	return
}

func (ms memStorage) GetCounterValue(name string) (value *int64, err error) {
	if currMetric, ok := ms.metricsMap[name]; ok {
		if currMetric.GetType() != "counter" {
			err = fmt.Errorf("metric type mismatch")
			return
		}
		value = currMetric.GetCounterValue()
	} else {
		err = fmt.Errorf("metric not found")
	}
	return
}

func (ms memStorage) ListAll() (metricsMap map[string]string, err error) {
	metricsMap = make(map[string]string, len(ms.metricsMap))
	for key, value := range ms.metricsMap {
		metricsMap[key] = value.GetValue()
	}
	return
}
