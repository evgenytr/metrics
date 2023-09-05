// Package config populates metric agent and server config data based on flags, environment variables or defaults
package config

import (
	"flag"
	"time"

	"github.com/evgenytr/metrics.git/internal/utils"

	"github.com/caarlos0/env/v6"
)

type agentConfig struct {
	Host           string  `env:"ADDRESS"`
	Key            string  `env:"KEY"`
	ReportInterval float64 `env:"REPORT_INTERVAL"`
	PollInterval   float64 `env:"POLL_INTERVAL"`
	RateLimit      int64   `env:"RATE_LIMIT"`
}

// GetAgentConfig returns agent config params
func GetAgentConfig() (host string, pollIntervalOut, reportIntervalOut time.Duration, key string, rateLimit int64) {

	var cfg agentConfig
	var pollIntervalIn, reportIntervalIn float64

	host, pollIntervalIn, reportIntervalIn, key, rateLimit = getAgentFlags()
	_ = env.Parse(&cfg)
	flag.Parse()

	if cfg.Host != "" {
		host = cfg.Host
	}

	if cfg.PollInterval != 0 {
		pollIntervalIn = cfg.PollInterval
	}

	if cfg.ReportInterval != 0 {
		reportIntervalIn = cfg.ReportInterval
	}

	if cfg.Key != "" {
		key = cfg.Key
	}

	if cfg.RateLimit != 0 {
		rateLimit = cfg.RateLimit
	}

	pollIntervalOut = utils.GetTimeInterval(pollIntervalIn)
	reportIntervalOut = utils.GetTimeInterval(reportIntervalIn)

	return
}

func getAgentFlags() (host string, pollInterval, reportInterval float64, key string, rateLimit int64) {
	if flag.Lookup("a") == nil {
		host = *flag.String("a", "localhost:8080", "host address")
	}
	if flag.Lookup("p") == nil {
		pollInterval = *flag.Float64("p", 2, "metrics polling interval")
	}
	if flag.Lookup("r") == nil {
		reportInterval = *flag.Float64("r", 10, "metrics reporting interval")
	}
	if flag.Lookup("k") == nil {
		key = *flag.String("k", "", "hash key")
	}
	if flag.Lookup("l") == nil {
		rateLimit = *flag.Int64("l", 2, "metrics report rate limit")
	}

	return
}
