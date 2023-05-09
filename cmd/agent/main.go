package main

import (
	"github.com/evgenytr/metrics.git/internal/monitor"
	"log"
	"time"
)

var host = "localhost"
var port = "8080"

type Monitor interface {
	PollMetrics() error
	ReportMetrics(host, port string) error
	ResetPollCount()
}

func main() {
	counter := 0
	pollInterval := 2 * time.Second
	//reportInterval := 10 * time.Second
	var currMetrics Monitor = monitor.GetNewMonitor()

	for {
		time.Sleep(pollInterval)
		counter++
		err := currMetrics.PollMetrics()
		if err != nil {
			log.Fatalln(err)
			panic(err)
		}
		if counter == 5 {
			counter = 0
			err = currMetrics.ReportMetrics(host, port)
			if err != nil {
				log.Fatalln(err)
				panic(err)
			}
			currMetrics.ResetPollCount()
		}

	}
}
