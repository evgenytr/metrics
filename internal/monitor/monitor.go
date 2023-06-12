package monitor

import (
	"fmt"
	"math/rand"
	"runtime"

	"github.com/go-resty/resty/v2"

	"github.com/evgenytr/metrics.git/internal/metric"
)

type monitor struct {
	Alloc,
	BuckHashSys,
	Frees,
	GCSys,
	HeapAlloc,
	HeapIdle,
	HeapInuse,
	HeapObjects,
	HeapReleased,
	HeapSys,
	LastGC,
	Lookups,
	MCacheInuse,
	MCacheSys,
	MSpanInuse,
	MSpanSys,
	Mallocs,
	NextGC,
	OtherSys,
	PauseTotalNs,
	StackInuse,
	StackSys,
	Sys,
	TotalAlloc,
	PollCount uint64

	NumForcedGC,
	NumGC uint32

	GCCPUFraction,
	RandomValue float64
}

type Monitor interface {
	PollMetrics() error
	ReportMetrics(host string) error
	ResetPollCount()
}

func NewMonitor() Monitor {
	return &monitor{}
}

func (m *monitor) ResetPollCount() {
	m.PollCount = 0
}

func (m *monitor) PollMetrics() (err error) {
	fmt.Println("pollMetrics")
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	m.Alloc = rtm.Alloc
	m.BuckHashSys = rtm.BuckHashSys
	m.Frees = rtm.Frees
	m.GCCPUFraction = rtm.GCCPUFraction
	m.GCSys = rtm.GCSys
	m.HeapAlloc = rtm.HeapAlloc
	m.HeapIdle = rtm.HeapIdle
	m.HeapInuse = rtm.HeapInuse
	m.HeapObjects = rtm.HeapObjects
	m.HeapReleased = rtm.HeapReleased
	m.HeapSys = rtm.HeapSys
	m.LastGC = rtm.LastGC
	m.Lookups = rtm.Lookups
	m.MCacheInuse = rtm.MCacheInuse
	m.MCacheSys = rtm.MCacheSys
	m.MSpanInuse = rtm.MSpanInuse
	m.MSpanSys = rtm.MSpanSys
	m.Mallocs = rtm.Mallocs
	m.NextGC = rtm.NextGC
	m.NumForcedGC = rtm.NumForcedGC
	m.NumGC = rtm.NumGC
	m.OtherSys = rtm.OtherSys
	m.PauseTotalNs = rtm.PauseTotalNs
	m.StackInuse = rtm.StackInuse
	m.StackSys = rtm.StackSys
	m.Sys = rtm.Sys
	m.TotalAlloc = rtm.TotalAlloc

	m.PollCount++
	m.RandomValue = rand.Float64()
	return
}

func (m *monitor) ReportMetrics(hostAddress string) (err error) {
	fmt.Println("reportMetrics")

	err = reportUint64Metric(metric.GaugeMetricType, "Alloc", m.Alloc, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "BuckHashSys", m.BuckHashSys, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "Frees", m.Frees, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "GCSys", m.GCSys, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "HeapAlloc", m.HeapAlloc, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "HeapIdle", m.HeapIdle, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "HeapInuse", m.HeapInuse, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "HeapObjects", m.HeapObjects, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "HeapReleased", m.HeapReleased, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "HeapSys", m.HeapSys, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "LastGC", m.LastGC, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "Lookups", m.Lookups, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "MCacheInuse", m.MCacheInuse, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "MCacheSys", m.MCacheSys, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "MSpanInuse", m.MSpanInuse, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "MSpanSys", m.MSpanSys, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "Mallocs", m.Mallocs, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "NextGC", m.NextGC, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "OtherSys", m.OtherSys, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "PauseTotalNs", m.PauseTotalNs, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "StackInuse", m.StackInuse, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "StackSys", m.StackSys, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "Sys", m.Sys, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "TotalAlloc", m.TotalAlloc, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.CounterMetricType, "PollCount", m.PollCount, hostAddress)
	if err != nil {
		return
	}

	err = reportUint64Metric(metric.GaugeMetricType, "NumForcedGC", uint64(m.NumForcedGC), hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric(metric.GaugeMetricType, "NumGC", uint64(m.NumGC), hostAddress)
	if err != nil {
		return
	}

	err = reportFloat64Metric(metric.GaugeMetricType, "GCCPUFraction", m.GCCPUFraction, hostAddress)
	if err != nil {
		return
	}
	err = reportFloat64Metric(metric.GaugeMetricType, "RandomValue", m.RandomValue, hostAddress)
	if err != nil {
		return
	}

	return
}

func reportFloat64Metric(metricType, name string, value float64, hostAddress string) (err error) {
	switch metricType {
	case metric.GaugeMetricType:
		err = reportGaugeMetric(name, value, hostAddress)
	case metric.CounterMetricType:
		err = reportCounterMetric(name, int64(value), hostAddress)
	default:
		err = fmt.Errorf("metric type not supported")
	}
	return
}

func reportUint64Metric(metricType, name string, value uint64, hostAddress string) (err error) {
	switch metricType {
	case metric.GaugeMetricType:
		err = reportGaugeMetric(name, float64(value), hostAddress)
	case metric.CounterMetricType:
		err = reportCounterMetric(name, int64(value), hostAddress)
	default:
		err = fmt.Errorf("metric type not supported")
	}
	return
}

func reportGaugeMetric(name string, value float64, hostAddress string) (err error) {
	currMetric := &metric.Metrics{ID: name, MType: metric.GaugeMetricType, Value: &value}
	err = postJSONMetric(currMetric, hostAddress)
	return
}

func reportCounterMetric(name string, value int64, hostAddress string) (err error) {
	currMetric := &metric.Metrics{ID: name, MType: metric.CounterMetricType, Delta: &value}
	err = postJSONMetric(currMetric, hostAddress)
	return
}

func postJSONMetric(metrics *metric.Metrics, hostAddress string) (err error) {

	client := resty.New()
	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(metrics).
		Post(fmt.Sprintf("%v/update/", hostAddress))

	//TODO: properly handle connection refused error (don't quit goroutine)
	//but quit on fatal error

	if err != nil {
		fmt.Println(err)
	}
	return nil
}
