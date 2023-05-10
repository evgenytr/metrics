package memstorage

import (
	"fmt"
	"github.com/evgenytr/metrics.git/internal/metric"
)

type MemStorage struct {
	metricsMap map[string]metric.Metric
}

func GetNewStorage() *MemStorage {
	return &MemStorage{
		metricsMap: make(map[string]metric.Metric),
	}
}

func (ms MemStorage) Update(metricType, name, value string) (err error) {
	if currMetric, ok := ms.metricsMap[name]; ok {
		err = currMetric.Add(metricType, value)
		if err != nil {
			return
		}
		//TODO get as pointer (?) instead of getting as value and reassigning
		ms.metricsMap[name] = currMetric
	} else {
		ms.metricsMap[name], err = metric.Create(metricType, name, value)
	}
	return
}

func (ms MemStorage) ReadValue(metricType, name string) (value string, err error) {
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
func (ms MemStorage) ListAll() (metricsMap map[string]string, err error) {
	metricsMap = make(map[string]string)
	for key, value := range ms.metricsMap {
		metricsMap[key] = value.GetValue()
	}
	return
}
