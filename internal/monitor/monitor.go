package monitor

import (
	"fmt"
	"math/rand"
	"runtime"

	"github.com/go-resty/resty/v2"

	"github.com/evgenytr/metrics.git/internal/metric"
)

type monitor struct {
	metrics map[string]*metric.Metrics
}

type Monitor interface {
	PollMetrics() error
	ReportMetrics(host string) error
	ResetPollCount()
}

func initMap() (initialMap map[string]*metric.Metrics, err error) {

	var metricList = map[string]string{
		"Alloc":         metric.GaugeMetricType,
		"BuckHashSys":   metric.GaugeMetricType,
		"Frees":         metric.GaugeMetricType,
		"GCSys":         metric.GaugeMetricType,
		"HeapAlloc":     metric.GaugeMetricType,
		"HeapIdle":      metric.GaugeMetricType,
		"HeapInuse":     metric.GaugeMetricType,
		"HeapObjects":   metric.GaugeMetricType,
		"HeapReleased":  metric.GaugeMetricType,
		"HeapSys":       metric.GaugeMetricType,
		"LastGC":        metric.GaugeMetricType,
		"Lookups":       metric.GaugeMetricType,
		"MCacheInuse":   metric.GaugeMetricType,
		"MCacheSys":     metric.GaugeMetricType,
		"MSpanInuse":    metric.GaugeMetricType,
		"MSpanSys":      metric.GaugeMetricType,
		"Mallocs":       metric.GaugeMetricType,
		"NextGC":        metric.GaugeMetricType,
		"OtherSys":      metric.GaugeMetricType,
		"PauseTotalNs":  metric.GaugeMetricType,
		"StackInuse":    metric.GaugeMetricType,
		"StackSys":      metric.GaugeMetricType,
		"Sys":           metric.GaugeMetricType,
		"TotalAlloc":    metric.GaugeMetricType,
		"PollCount":     metric.CounterMetricType,
		"NumForcedGC":   metric.GaugeMetricType,
		"NumGC":         metric.GaugeMetricType,
		"GCCPUFraction": metric.GaugeMetricType,
		"RandomValue":   metric.GaugeMetricType,
	}

	initialMap = make(map[string]*metric.Metrics, len(metricList))
	for metricName, metricType := range metricList {
		switch metricType {
		case metric.GaugeMetricType:
			var value float64
			initialMap[metricName], err = metric.CreateGauge(metricName, &value)
		case metric.CounterMetricType:
			var value int64
			initialMap[metricName], err = metric.CreateCounter(metricName, &value)
		}
		if err != nil {
			return
		}
	}

	return
}

func NewMonitor() (m Monitor, err error) {
	initialMap, err := initMap()
	if err != nil {
		return
	}
	m = &monitor{
		metrics: initialMap,
	}
	return
}

func (m *monitor) ResetPollCount() {
	_ = m.metrics["PollCount"].ResetCounter()
}

func updateGaugeMetric(metric *metric.Metrics, value float64) (err error) {
	_, err = metric.UpdateGauge(&value)
	return
}

func updateCounterMetric(metric *metric.Metrics, value int64) (err error) {
	_, err = metric.UpdateCounter(&value)
	return
}

func (m *monitor) PollMetrics() (err error) {
	fmt.Println("pollMetrics")
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	err = updateGaugeMetric(m.metrics["Alloc"], float64(rtm.Alloc))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["BuckHashSys"], float64(rtm.BuckHashSys))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["Frees"], float64(rtm.Frees))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["GCCPUFraction"], rtm.GCCPUFraction)
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["GCSys"], float64(rtm.GCSys))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["HeapAlloc"], float64(rtm.HeapAlloc))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["HeapIdle"], float64(rtm.HeapIdle))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["HeapInuse"], float64(rtm.HeapInuse))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["HeapObjects"], float64(rtm.HeapObjects))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["HeapReleased"], float64(rtm.HeapReleased))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["HeapSys"], float64(rtm.HeapSys))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["LastGC"], float64(rtm.LastGC))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["Lookups"], float64(rtm.Lookups))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["MCacheInuse"], float64(rtm.MCacheInuse))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["MCacheSys"], float64(rtm.MCacheSys))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["MSpanInuse"], float64(rtm.MSpanInuse))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["MSpanSys"], float64(rtm.MSpanSys))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["NextGC"], float64(rtm.NextGC))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["NumForcedGC"], float64(rtm.NumForcedGC))
	if err != nil {
		return
	}
	err = updateGaugeMetric(m.metrics["NumGC"], float64(rtm.NumGC))
	if err != nil {
		return
	}
	err = updateGaugeMetric(m.metrics["OtherSys"], float64(rtm.OtherSys))
	if err != nil {
		return
	}
	err = updateGaugeMetric(m.metrics["PauseTotalNs"], float64(rtm.PauseTotalNs))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["StackInuse"], float64(rtm.StackInuse))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["StackSys"], float64(rtm.StackSys))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["Sys"], float64(rtm.Sys))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["Mallocs"], float64(rtm.Mallocs))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["TotalAlloc"], float64(rtm.TotalAlloc))
	if err != nil {
		return
	}

	err = updateCounterMetric(m.metrics["PollCount"], 1)
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["RandomValue"], rand.Float64())
	if err != nil {
		return
	}

	return
}

func (m *monitor) ReportMetrics(hostAddress string) (err error) {
	fmt.Println("reportMetrics")

	if len(m.metrics) == 0 {
		fmt.Println("empty batch")
		return
	}

	var metricsBatch []metric.Metrics

	for _, value := range m.metrics {
		metricsBatch = append(metricsBatch, *value)
	}

	client := resty.New()
	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(metricsBatch).
		Post(fmt.Sprintf("%v/updates/", hostAddress))

	//TODO: properly handle connection refused error (don't quit goroutine)
	//but quit on fatal error

	if err != nil {
		fmt.Println(err)
	}
	return nil
}
