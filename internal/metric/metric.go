package metric

import (
	"fmt"
	"strconv"
)

// TODO having 2 structs is redundant, keep only one
type Metric struct {
	Name       string  `json:"id"`
	MetricType string  `json:"type"`
	Gauge      float64 `json:"gauge,omitempty"`
	Counter    int64   `json:"counter,omitempty"`
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func Create(metricType, name, value string) (newMetric *Metric, err error) {
	fmt.Println("Metric Create")
	newMetric = &Metric{
		Name:       name,
		MetricType: metricType,
	}
	_, err = newMetric.Add(metricType, value)
	return
}

func CreateGauge(name string, value *float64) (newMetric *Metric, err error) {
	fmt.Println("Gauge Metric Create")
	newMetric = &Metric{
		Name:       name,
		MetricType: "gauge",
		Gauge:      *value,
	}
	return
}

func CreateCounter(name string, value *int64) (newMetric *Metric, err error) {
	fmt.Println("Counter Metric Create")
	newMetric = &Metric{
		Name:       name,
		MetricType: "counter",
		Counter:    *value,
	}
	return
}

func (metric *Metric) GetValue() (value string) {

	switch metric.MetricType {
	case "gauge":
		value = strconv.FormatFloat(metric.Gauge, 'f', -1, 64)
	case "counter":
		value = strconv.FormatInt(metric.Counter, 10)
	}
	return
}

func (metric *Metric) GetGaugeValue() (value *float64) {
	return &metric.Gauge
}

func (metric *Metric) GetCounterValue() (value *int64) {
	return &metric.Counter
}

func (metric *Metric) GetType() (value string) {
	return metric.MetricType
}

func (metric *Metric) Add(metricType, value string) (newValue string, err error) {
	fmt.Println("Metric Add")
	if metric.MetricType != metricType {
		err = fmt.Errorf("metric type mismatch")
		return
	}

	switch metric.MetricType {
	case "gauge":
		var floatValue float64
		floatValue, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return
		}
		if floatValue < 0 {
			err = fmt.Errorf("value less than zero")
			return
		}
		metric.Gauge = floatValue
		newValue = strconv.FormatFloat(floatValue, 'f', -1, 64)
	case "counter":
		var intValue int64
		intValue, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return
		}
		if intValue < 0 {
			err = fmt.Errorf("value less than zero")
			return
		}
		metric.Counter += intValue
		newValue = strconv.FormatInt(metric.Counter, 10)
	default:
		err = fmt.Errorf("metric type not supported")
	}
	return
}

func (metric *Metric) UpdateGauge(value *float64) (newValue *float64, err error) {

	if metric.MetricType != "gauge" {
		err = fmt.Errorf("metric type mismatch")
		return
	}

	metric.Gauge = *value
	newValue = value
	return
}

func (metric *Metric) UpdateCounter(value *int64) (newValue *int64, err error) {

	if metric.MetricType != "counter" {
		err = fmt.Errorf("metric type mismatch")
		return
	}

	metric.Counter += *value
	newValue = &metric.Counter
	return
}
