package metric

import (
	"fmt"
	"strconv"
)

const GaugeMetricType = "gauge"
const CounterMetricType = "counter"

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func Create(metricType, name, value string) (newMetric *Metrics, err error) {
	fmt.Println("Metric Create")
	newMetric = &Metrics{
		ID:    name,
		MType: metricType,
	}
	_, err = newMetric.Add(metricType, value)
	return
}

func CreateGauge(name string, value *float64) (newMetric *Metrics, err error) {
	fmt.Println("Gauge Metric Create")
	newMetric = &Metrics{
		ID:    name,
		MType: GaugeMetricType,
		Value: value,
	}
	return
}

func CreateCounter(name string, value *int64) (newMetric *Metrics, err error) {
	fmt.Println("Counter Metric Create")
	newMetric = &Metrics{
		ID:    name,
		MType: CounterMetricType,
		Delta: value,
	}
	return
}

func (metric *Metrics) GetValue() (value string) {

	switch metric.MType {
	case GaugeMetricType:
		value = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
	case CounterMetricType:
		value = strconv.FormatInt(*metric.Delta, 10)
	}
	return
}

func (metric *Metrics) GetGaugeValue() (value *float64) {
	return metric.Value
}

func (metric *Metrics) GetCounterValue() (value *int64) {
	return metric.Delta
}

func (metric *Metrics) GetType() (value string) {
	return metric.MType
}

func (metric *Metrics) Add(metricType, value string) (newValue string, err error) {
	fmt.Println("Metric Add")
	if metric.MType != metricType {
		err = fmt.Errorf("metric type mismatch")
		return
	}

	switch metric.MType {
	case GaugeMetricType:
		var floatValue float64
		floatValue, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return
		}
		if floatValue < 0 {
			err = fmt.Errorf("value less than zero")
			return
		}
		metric.Value = &floatValue
		newValue = strconv.FormatFloat(floatValue, 'f', -1, 64)
	case CounterMetricType:
		var intValue int64
		intValue, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return
		}
		if intValue < 0 {
			err = fmt.Errorf("value less than zero")
			return
		}
		if metric.Delta != nil {
			intValue += *metric.Delta
		}
		metric.Delta = &intValue
		newValue = strconv.FormatInt(intValue, 10)
	default:
		err = fmt.Errorf("metric type not supported")
	}
	return
}

func (metric *Metrics) UpdateGauge(value *float64) (newValue *float64, err error) {

	if metric.MType != GaugeMetricType {
		err = fmt.Errorf("metric type mismatch")
		return
	}

	metric.Value = value
	newValue = value
	return
}

func (metric *Metrics) UpdateCounter(value *int64) (newValue *int64, err error) {

	if metric.MType != CounterMetricType {
		err = fmt.Errorf("metric type mismatch")
		return
	}

	deltaValue := *value
	if metric.Delta != nil {
		deltaValue += *metric.Delta
	}
	metric.Delta = &deltaValue
	newValue = metric.Delta
	return
}
