package metric

import (
	"fmt"
	"strconv"
)

type Metric struct {
	name       string
	metricType string
	gauge      float64
	counter    int64
}

func Create(metricType, name, value string) (newMetric *Metric, err error) {
	fmt.Println("Metric Create")
	newMetric = &Metric{
		name:       name,
		metricType: metricType,
	}
	err = newMetric.Add(metricType, value)
	return
}

func (metric *Metric) GetValue() (value string) {

	switch metric.metricType {
	case "gauge":
		value = strconv.FormatFloat(metric.gauge, 'f', 3, 64)
	case "counter":
		value = strconv.FormatInt(metric.counter, 10)
	}
	return
}

func (metric *Metric) GetType() (value string) {
	return metric.metricType
}

func (metric *Metric) Add(metricType, value string) error {
	fmt.Println("Metric Add")
	if metric.metricType != metricType {
		return fmt.Errorf("metric type mismatch")
	}

	switch metric.metricType {
	case "gauge":
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("cannot parse float")
		}
		if floatValue < 0 {
			return fmt.Errorf("value less than zero")
		}
		metric.gauge = floatValue
	case "counter":
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot parse int")
		}
		if intValue < 0 {
			return fmt.Errorf("value less than zero")
		}
		metric.counter += intValue
	default:
		return fmt.Errorf("metric type not supported")
	}
	return nil
}
