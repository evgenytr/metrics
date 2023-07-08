package main

import (
	"context"
	"fmt"
	"github.com/evgenytr/metrics.git/internal/config"
	"github.com/evgenytr/metrics.git/internal/monitor"
	"log"
)

func main() {

	host, pollInterval, reportInterval, key, rateLimit := config.GetAgentConfig()

	fmt.Println(*host, *pollInterval, *reportInterval, *key, *rateLimit)
	hostAddress := fmt.Sprintf("http://%v", *host)

	currMetrics, err := monitor.NewMonitor(&hostAddress, key)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	workerCtx, cancelWorkerCtx := context.WithCancelCause(ctx)
	//create poll and report queues
	pollQueue := monitor.NewQueue(nil)
	extraPollQueue := monitor.NewQueue(nil)
	reportQueue := monitor.NewQueue(rateLimit)

	//create workers
	var pollWorkerID int64 = 0
	var extraPollWorkerID int64 = 1
	var reportWorkerID int64 = 2

	pollWorker := monitor.NewWorker(pollWorkerID, pollQueue)
	go pollWorker.Loop(workerCtx, cancelWorkerCtx, currMetrics.PollMetrics)

	extraPollWorker := monitor.NewWorker(extraPollWorkerID, extraPollQueue)
	go extraPollWorker.Loop(workerCtx, cancelWorkerCtx, currMetrics.PollAdditionalMetrics)

	if rateLimit == nil || *rateLimit <= 0 {
		reportWorker := monitor.NewWorker(reportWorkerID, reportQueue)
		go reportWorker.Loop(workerCtx, cancelWorkerCtx, currMetrics.ReportMetrics)
	} else {
		for i := int64(0); i < *rateLimit; i++ {
			reportWorker := monitor.NewWorker(i+reportWorkerID, reportQueue)
			go reportWorker.Loop(workerCtx, cancelWorkerCtx, currMetrics.ReportMetrics)
		}
	}

	//fill queues with tasks according to time intervals
	go pollQueue.ScheduleTasks(pollInterval)
	go extraPollQueue.ScheduleTasks(pollInterval)
	go reportQueue.ScheduleTasks(reportInterval)
	
	for {
		<-ctx.Done()
		err := context.Cause(ctx)
		if err != nil {
			log.Fatalln(err)
		}
	}

}
