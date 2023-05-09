package memstorage

import (
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
	} else {
		ms.metricsMap[name], err = metric.Create(metricType, name, value)
	}
	return
}
