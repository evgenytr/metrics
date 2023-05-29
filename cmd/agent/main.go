package main

import (
	"context"
	"fmt"
	"github.com/evgenytr/metrics.git/internal/config"
	"github.com/evgenytr/metrics.git/internal/monitor"
	"log"
	"time"
)

func main() {

	host, pollInterval, reportInterval := config.GetAgentConfig()

	fmt.Println(*host, *pollInterval, *reportInterval)

	var currMetrics = monitor.NewMonitor()

	hostAddress := fmt.Sprintf("http://%v", *host)

	ctx := context.Background()
	go pollMetrics(ctx, *pollInterval, currMetrics)
	go reportMetrics(ctx, *reportInterval, currMetrics, hostAddress)

	for {
		<-ctx.Done()
		err := context.Cause(ctx)
		if err != nil {
			log.Fatalln(err)
		}
	}

}

func pollMetrics(ctx context.Context, pollInterval float64, currMetrics monitor.Monitor) {
	_, cancelCtx := context.WithCancelCause(ctx)
	for {
		time.Sleep(time.Duration(pollInterval) * time.Second)
		err := currMetrics.PollMetrics()

		if err != nil {
			fmt.Println("poll metrics err")
			cancelCtx(err)
			return
		}
	}
}

func reportMetrics(ctx context.Context, reportInterval float64, currMetrics monitor.Monitor, host string) {
	_, cancelCtx := context.WithCancelCause(ctx)
	for {
		time.Sleep(time.Duration(reportInterval) * time.Second)
		err := currMetrics.ReportMetrics(host)

		if err != nil {
			cancelCtx(err)
			return
		}

		currMetrics.ResetPollCount()
	}
}
