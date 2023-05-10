package main

import (
	"flag"
	"fmt"
	"github.com/evgenytr/metrics.git/internal/monitor"
	"log"
	"time"
)

type Monitor interface {
	PollMetrics() error
	ReportMetrics(host string) error
	ResetPollCount()
}

var (
	host           *string
	pollInterval   *float64
	reportInterval *float64
)

func init() {
	host = flag.String("a", "localhost:8080", "host address")
	pollInterval = flag.Float64("p", 2, "metrics polling interval")
	reportInterval = flag.Float64("r", 10, "metrics reporting interval")
}

func main() {
	flag.Parse()

	var currMetrics Monitor = monitor.GetNewMonitor()
	var errChannel = make(chan error)
	hostAddress := fmt.Sprintf("http://%v", *host)
	go pollMetrics(*pollInterval, currMetrics, errChannel)
	go reportMetrics(*reportInterval, currMetrics, hostAddress, errChannel)

	err := <-errChannel

	if err != nil {
		log.Fatalln(err)
	}
}

func pollMetrics(pollInterval float64, currMetrics Monitor, errChannel chan error) {
	for {
		time.Sleep(time.Duration(pollInterval) * time.Second)
		err := currMetrics.PollMetrics()

		if err != nil {
			errChannel <- err
			return
		}
	}
}

func reportMetrics(reportInterval float64, currMetrics Monitor, host string, errChannel chan error) {
	for {
		time.Sleep(time.Duration(reportInterval) * time.Second)
		err := currMetrics.ReportMetrics(host)

		if err != nil {
			errChannel <- err
			return
		}
		currMetrics.ResetPollCount()
	}
}
