package monitor

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"runtime"
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

	err = reportUint64Metric("gauge", "Alloc", m.Alloc, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "BuckHashSys", m.BuckHashSys, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "Frees", m.Frees, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "GCSys", m.GCSys, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "HeapAlloc", m.HeapAlloc, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "HeapIdle", m.HeapIdle, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "HeapInuse", m.HeapInuse, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "HeapObjects", m.HeapObjects, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "HeapReleased", m.HeapReleased, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "HeapSys", m.HeapSys, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "LastGC", m.LastGC, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "Lookups", m.Lookups, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "MCacheInuse", m.MCacheInuse, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "MCacheSys", m.MCacheSys, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "MSpanInuse", m.MSpanInuse, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "MSpanSys", m.MSpanSys, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "Mallocs", m.Mallocs, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "NextGC", m.NextGC, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "OtherSys", m.OtherSys, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "PauseTotalNs", m.PauseTotalNs, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "StackInuse", m.StackInuse, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "StackSys", m.StackSys, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "Sys", m.Sys, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("gauge", "TotalAlloc", m.TotalAlloc, hostAddress)
	if err != nil {
		return
	}
	err = reportUint64Metric("counter", "PollCount", m.PollCount, hostAddress)
	if err != nil {
		return
	}

	err = reportUint32Metric("gauge", "NumForcedGC", m.NumForcedGC, hostAddress)
	if err != nil {
		return
	}
	err = reportUint32Metric("gauge", "NumGC", m.NumGC, hostAddress)
	if err != nil {
		return
	}

	err = reportFloat64Metric("gauge", "GCCPUFraction", m.GCCPUFraction, hostAddress)
	if err != nil {
		return
	}
	err = reportFloat64Metric("gauge", "RandomValue", m.RandomValue, hostAddress)
	if err != nil {
		return
	}

	return
}

func reportUint64Metric(metricType, name string, value uint64, hostAddress string) (err error) {

	client := resty.New()
	_, err = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%v/update/%v/%v/%v", hostAddress, metricType, name, value))

	return
}

func reportFloat64Metric(metricType, name string, value float64, hostAddress string) (err error) {
	client := resty.New()
	_, err = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%v/update/%v/%v/%v", hostAddress, metricType, name, value))
	return
}

func reportUint32Metric(metricType, name string, value uint32, hostAddress string) (err error) {
	client := resty.New()
	_, err = client.R().
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%v/update/%v/%v/%v", hostAddress, metricType, name, value))
	return
}
