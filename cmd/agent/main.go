package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/evgenytr/metrics.git/internal/config"
	errorHandling "github.com/evgenytr/metrics.git/internal/errors"
	"github.com/evgenytr/metrics.git/internal/monitor"
)

func main() {

	host, pollInterval, reportInterval, key, rateLimit := config.GetAgentConfig()

	fmt.Println(*host, *pollInterval, *reportInterval, *key, *rateLimit)

	currMetrics, err := monitor.NewMonitor()
	if err != nil {
		log.Fatalln(err)
	}

	hostAddress := fmt.Sprintf("http://%v", *host)

	ctx := context.Background()
	go pollMetrics(ctx, pollInterval, currMetrics)
	go reportMetrics(ctx, reportInterval, currMetrics, &hostAddress, key)
	go pollAdditionalMetrics(ctx, pollInterval, currMetrics)

	for {
		<-ctx.Done()
		err := context.Cause(ctx)
		if err != nil {
			log.Fatalln(err)
		}
	}

}

func pollMetrics(ctx context.Context, pollInterval *time.Duration, currMetrics monitor.Monitor) {
	_, cancelCtx := context.WithCancelCause(ctx)
	for {
		time.Sleep(*pollInterval)
		err := currMetrics.PollMetrics()

		if err != nil {
			fmt.Println("poll metrics err")
			cancelCtx(err)
			return
		}
	}
}

func pollAdditionalMetrics(ctx context.Context, pollInterval *time.Duration, currMetrics monitor.Monitor) {
	_, cancelCtx := context.WithCancelCause(ctx)
	for {
		time.Sleep(*pollInterval)
		err := currMetrics.PollAdditionalMetrics()

		if err != nil {
			fmt.Println("poll additional metrics err")
			cancelCtx(err)
			return
		}
	}
}

func reportMetrics(ctx context.Context, reportInterval *time.Duration, currMetrics monitor.Monitor, host, key *string) {
	_, cancelCtx := context.WithCancelCause(ctx)
	for {
		time.Sleep(*reportInterval)
		err := currMetrics.ReportMetrics(host, key)

		if err != nil {
			for _, retryInterval := range errorHandling.RepeatedAttemptsIntervals {
				time.Sleep(*retryInterval)
				err = currMetrics.ReportMetrics(host, key)
				if err == nil {
					break
				}
			}
		}

		if err != nil {
			cancelCtx(err)
			return
		}

		currMetrics.ResetPollCount()
	}
}
