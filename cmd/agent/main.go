// Package main initializes metrics agent service and starts it.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/evgenytr/metrics.git/internal/config"
	"github.com/evgenytr/metrics.git/internal/monitor"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {

	fmt.Println("Build version: ", buildVersion)
	fmt.Println("Build date: ", buildDate)
	fmt.Println("Build commit: ", buildCommit)

	host, hostGrpc, pollInterval, reportInterval, key, rateLimit, cryptoKey := config.GetAgentConfig()

	fmt.Println(host, pollInterval, reportInterval, key, rateLimit, cryptoKey)
	hostAddress := fmt.Sprintf("http://%v", host)

	var wg sync.WaitGroup
	currMetrics, err := monitor.NewMonitor(hostAddress, hostGrpc, key, cryptoKey, &wg)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	workerCtx, cancelWorkerCtx := context.WithCancelCause(ctx)
	queueCtx, stopQueueCtx := context.WithCancelCause(ctx)

	//create poll and report queues
	pollQueue := monitor.NewQueue(0)
	extraPollQueue := monitor.NewQueue(0)
	reportQueue := monitor.NewQueue(rateLimit)

	//create workers
	var pollWorkerID int64 = 0
	var extraPollWorkerID int64 = 1
	var reportWorkerID int64 = 2

	pollWorker := monitor.NewWorker(pollWorkerID, pollQueue)
	go pollWorker.Loop(workerCtx, cancelWorkerCtx, currMetrics.PollMetrics)

	extraPollWorker := monitor.NewWorker(extraPollWorkerID, extraPollQueue)
	go extraPollWorker.Loop(workerCtx, cancelWorkerCtx, currMetrics.PollAdditionalMetrics)

	if rateLimit <= 0 {
		reportWorker := monitor.NewWorker(reportWorkerID, reportQueue)
		go reportWorker.Loop(workerCtx, cancelWorkerCtx, currMetrics.ReportMetricsGrpc)
	} else {
		for i := int64(0); i < rateLimit; i++ {
			reportWorker := monitor.NewWorker(i+reportWorkerID, reportQueue)
			go reportWorker.Loop(workerCtx, cancelWorkerCtx, currMetrics.ReportMetricsGrpc)
		}
	}

	//fill queues with tasks according to time intervals
	go pollQueue.ScheduleTasks(queueCtx, pollInterval)
	go extraPollQueue.ScheduleTasks(queueCtx, pollInterval)
	go reportQueue.ScheduleTasks(queueCtx, reportInterval)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	select {
	case <-workerCtx.Done():
		err := context.Cause(workerCtx)
		if err != nil {
			log.Fatalln(err)
		}

	case <-sigChan:
		fmt.Println("shutdown signal")
		err = fmt.Errorf("shutdown signal received")
		stopQueueCtx(err)
		cancelWorkerCtx(err)

		wg.Wait()
		fmt.Println("after wait group done")

		fmt.Println("report collected metrics before shutdown")
		//report collected metrics before shutting down
		err = currMetrics.ReportMetricsGrpc()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("metrics sent")

	}

}
