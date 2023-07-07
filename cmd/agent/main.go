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

	//create poll and report queues
	pollQueue := monitor.NewQueue(nil)
	extraPollQueue := monitor.NewQueue(nil)
	reportQueue := monitor.NewQueue(rateLimit)

	//create workers
	var pollWorkerId int64 = 0
	var extraPollWorkerId int64 = 1
	var reportWorkerId int64 = 2

	pollWorker := monitor.NewWorker(pollWorkerId, pollQueue)
	go pollWorker.Loop(ctx, currMetrics.PollMetrics)

	extraPollWorker := monitor.NewWorker(extraPollWorkerId, extraPollQueue)
	go extraPollWorker.Loop(ctx, currMetrics.PollAdditionalMetrics)

	if rateLimit == nil || *rateLimit <= 0 {
		reportWorker := monitor.NewWorker(reportWorkerId, reportQueue)
		go reportWorker.Loop(ctx, currMetrics.ReportMetrics)
	} else {
		for i := reportWorkerId; i < *rateLimit; i++ {
			reportWorker := monitor.NewWorker(i, reportQueue)
			go reportWorker.Loop(ctx, currMetrics.ReportMetrics)
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
