package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/evgenytr/metrics.git/internal/config"
	"github.com/evgenytr/metrics.git/internal/monitor"
	"log"
	"time"
)

func main() {

	host, pollInterval, reportInterval := getFlags()
	var cfg config.Config
	_ = env.Parse(&cfg)
	flag.Parse()

	if cfg.Host != "" {
		host = &cfg.Host
	}

	if cfg.PollInterval != 0 {
		pollInterval = &(cfg.PollInterval)
	}

	if cfg.ReportInterval != 0 {
		reportInterval = &(cfg.ReportInterval)
	}

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

func getFlags() (host *string, pollInterval, reportInterval *float64) {
	host = flag.String("a", "localhost:8080", "host address")
	pollInterval = flag.Float64("p", 2, "metrics polling interval")
	reportInterval = flag.Float64("r", 10, "metrics reporting interval")
	return
}
