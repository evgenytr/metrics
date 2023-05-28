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

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func Create(metricType, name, value string) (newMetric *Metric, err error) {
	fmt.Println("Metric Create")
	newMetric = &Metric{
		name:       name,
		metricType: metricType,
	}
	_, err = newMetric.Add(metricType, value)
	return
}

func CreateGauge(name string, value *float64) (newMetric *Metric, err error) {
	fmt.Println("Gauge Metric Create")
	newMetric = &Metric{
		name:       name,
		metricType: "gauge",
		gauge:      *value,
	}
	return
}

func CreateCounter(name string, value *int64) (newMetric *Metric, err error) {
	fmt.Println("Counter Metric Create")
	newMetric = &Metric{
		name:       name,
		metricType: "counter",
		counter:    *value,
	}
	return
}

func (metric *Metric) GetValue() (value string) {

	switch metric.metricType {
	case "gauge":
		value = strconv.FormatFloat(metric.gauge, 'f', -1, 64)
	case "counter":
		value = strconv.FormatInt(metric.counter, 10)
	}
	return
}

func (metric *Metric) GetGaugeValue() (value *float64) {
	return &metric.gauge
}

func (metric *Metric) GetCounterValue() (value *int64) {
	return &metric.counter
}

func (metric *Metric) GetType() (value string) {
	return metric.metricType
}

func (metric *Metric) Add(metricType, value string) (newValue string, err error) {
	fmt.Println("Metric Add")
	if metric.metricType != metricType {
		err = fmt.Errorf("metric type mismatch")
		return
	}

	switch metric.metricType {
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
		metric.gauge = floatValue
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
		metric.counter += intValue
		newValue = strconv.FormatInt(metric.counter, 10)
	default:
		err = fmt.Errorf("metric type not supported")
	}
	return
}

func (metric *Metric) UpdateGauge(value *float64) (newValue *float64, err error) {

	if metric.metricType != "gauge" {
		err = fmt.Errorf("metric type mismatch")
		return
	}

	metric.gauge = *value
	newValue = value
	return
}

func (metric *Metric) UpdateCounter(value *int64) (newValue *int64, err error) {

	if metric.metricType != "counter" {
		err = fmt.Errorf("metric type mismatch")
		return
	}

	metric.counter += *value
	newValue = &metric.counter
	return
}
