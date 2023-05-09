package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type Monitor struct {
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

var host = "localhost"
var port = "8080"

func main() {
	counter := 0
	pollInterval := 2 * time.Second
	//reportInterval := 10 * time.Second
	var currMetrics Monitor

	for {
		time.Sleep(pollInterval)
		counter++
		err := pollMetrics(&currMetrics)
		if err != nil {
			panic(err)
		}
		if counter == 4 {
			counter = 0
			err = reportMetrics(&currMetrics, host, port)
			if err != nil {
				panic(err)
			}
		}

	}
}

func pollMetrics(monitor *Monitor) (err error) {
	fmt.Println("pollMetrics")
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	monitor.Alloc = rtm.Alloc
	monitor.BuckHashSys = rtm.BuckHashSys
	monitor.Frees = rtm.Frees
	monitor.GCCPUFraction = rtm.GCCPUFraction
	monitor.GCSys = rtm.GCSys
	monitor.HeapAlloc = rtm.HeapAlloc
	monitor.HeapIdle = rtm.HeapIdle
	monitor.HeapInuse = rtm.HeapInuse
	monitor.HeapObjects = rtm.HeapObjects
	monitor.HeapReleased = rtm.HeapReleased
	monitor.HeapSys = rtm.HeapSys
	monitor.LastGC = rtm.LastGC
	monitor.Lookups = rtm.Lookups
	monitor.MCacheInuse = rtm.MCacheInuse
	monitor.MCacheSys = rtm.MCacheSys
	monitor.MSpanInuse = rtm.MSpanInuse
	monitor.MSpanSys = rtm.MSpanSys
	monitor.Mallocs = rtm.Mallocs
	monitor.NextGC = rtm.NextGC
	monitor.NumForcedGC = rtm.NumForcedGC
	monitor.NumGC = rtm.NumGC
	monitor.OtherSys = rtm.OtherSys
	monitor.PauseTotalNs = rtm.PauseTotalNs
	monitor.StackInuse = rtm.StackInuse
	monitor.StackSys = rtm.StackSys
	monitor.Sys = rtm.Sys
	monitor.TotalAlloc = rtm.TotalAlloc

	monitor.PollCount++
	monitor.RandomValue = rand.Float64()
	return
}

func reportMetrics(monitor *Monitor, host, port string) (err error) {
	fmt.Println("reportMetrics")
	var hostAddress = fmt.Sprintf("http://%v:%v", host, port)
	reportUint64Metric("gauge", "Alloc", monitor.Alloc, hostAddress)
	reportUint64Metric("gauge", "BuckHashSys", monitor.BuckHashSys, hostAddress)
	reportUint64Metric("gauge", "Frees", monitor.Frees, hostAddress)
	reportUint64Metric("gauge", "GCSys", monitor.GCSys, hostAddress)
	reportUint64Metric("gauge", "HeapAlloc", monitor.HeapAlloc, hostAddress)
	reportUint64Metric("gauge", "HeapIdle", monitor.HeapIdle, hostAddress)
	reportUint64Metric("gauge", "HeapInuse", monitor.HeapInuse, hostAddress)
	reportUint64Metric("gauge", "HeapObjects", monitor.HeapObjects, hostAddress)
	reportUint64Metric("gauge", "HeapReleased", monitor.HeapReleased, hostAddress)
	reportUint64Metric("gauge", "HeapSys", monitor.HeapSys, hostAddress)
	reportUint64Metric("gauge", "LastGC", monitor.LastGC, hostAddress)
	reportUint64Metric("gauge", "Lookups", monitor.Lookups, hostAddress)
	reportUint64Metric("gauge", "MCacheInuse", monitor.MCacheInuse, hostAddress)
	reportUint64Metric("gauge", "MCacheSys", monitor.MCacheSys, hostAddress)
	reportUint64Metric("gauge", "MSpanInuse", monitor.MSpanInuse, hostAddress)
	reportUint64Metric("gauge", "MSpanSys", monitor.MSpanSys, hostAddress)
	reportUint64Metric("gauge", "Mallocs", monitor.Mallocs, hostAddress)
	reportUint64Metric("gauge", "NextGC", monitor.NextGC, hostAddress)
	reportUint64Metric("gauge", "OtherSys", monitor.OtherSys, hostAddress)
	reportUint64Metric("gauge", "PauseTotalNs", monitor.PauseTotalNs, hostAddress)
	reportUint64Metric("gauge", "StackInuse", monitor.StackInuse, hostAddress)
	reportUint64Metric("gauge", "StackSys", monitor.StackSys, hostAddress)
	reportUint64Metric("gauge", "Sys", monitor.Sys, hostAddress)
	reportUint64Metric("gauge", "TotalAlloc", monitor.TotalAlloc, hostAddress)
	reportUint64Metric("counter", "PollCount", monitor.PollCount, hostAddress)
	monitor.PollCount = 0

	reportUint32Metric("gauge", "NumForcedGC", monitor.NumForcedGC, hostAddress)
	reportUint32Metric("gauge", "NumGC", monitor.NumGC, hostAddress)

	reportFloat64Metric("gauge", "GCCPUFraction", monitor.GCCPUFraction, hostAddress)
	reportFloat64Metric("gauge", "RandomValue", monitor.RandomValue, hostAddress)

	return
}

func reportUint64Metric(metricType, name string, value uint64, hostAddress string) (err error) {
	_, err = http.Post(fmt.Sprintf("%v/update/%v/%v/%v", hostAddress, metricType, name, value), "text/plain", nil)
	return
}

func reportFloat64Metric(metricType, name string, value float64, hostAddress string) (err error) {
	_, err = http.Post(fmt.Sprintf("%v/update/%v/%v/%v", hostAddress, metricType, name, value), "text/plain", nil)
	return
}

func reportUint32Metric(metricType, name string, value uint32, hostAddress string) (err error) {
	_, err = http.Post(fmt.Sprintf("%v/update/%v/%v/%v", hostAddress, metricType, name, value), "text/plain", nil)
	return
}
