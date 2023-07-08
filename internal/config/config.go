package config

import (
	"flag"
	"os"
	"time"

	"github.com/evgenytr/metrics.git/internal/utils"

	"github.com/caarlos0/env/v6"
)

type agentConfig struct {
	Host           string  `env:"ADDRESS"`
	ReportInterval float64 `env:"REPORT_INTERVAL"`
	PollInterval   float64 `env:"POLL_INTERVAL"`
	Key            string  `env:"KEY"`
	RateLimit      int64   `env:"RATE_LIMIT"`
}

type serverConfig struct {
	Host            string  `env:"ADDRESS"`
	StoreInterval   float64 `env:"STORE_INTERVAL"`
	FileStoragePath string  `env:"FILE_STORAGE_PATH"`
	Restore         bool    `env:"RESTORE"`
	DatabaseDSN     string  `env:"DATABASE_DSN"`
	Key             string  `env:"KEY"`
}

func GetAgentConfig() (host *string, pollIntervalOut, reportIntervalOut *time.Duration, key *string, rateLimit *int64) {

	var cfg agentConfig
	var pollIntervalIn, reportIntervalIn *float64

	host, pollIntervalIn, reportIntervalIn, key, rateLimit = getAgentFlags()
	_ = env.Parse(&cfg)
	flag.Parse()

	if cfg.Host != "" {
		host = &cfg.Host
	}

	if cfg.PollInterval != 0 {
		pollIntervalIn = &cfg.PollInterval
	}

	if cfg.ReportInterval != 0 {
		reportIntervalIn = &cfg.ReportInterval
	}

	if cfg.Key != "" {
		key = &cfg.Key
	}

	if cfg.RateLimit != 0 {
		rateLimit = &cfg.RateLimit
	}

	pollIntervalOut = utils.GetTimeInterval(*pollIntervalIn)
	reportIntervalOut = utils.GetTimeInterval(*reportIntervalIn)

	return
}
func GetServerConfig() (host *string, storeIntervalOut *time.Duration, fileStoragePath *string, restore *bool, dbDSN, key *string) {

	var storeIntervalIn *float64
	var cfg serverConfig

	host, storeIntervalIn, fileStoragePath, restore, dbDSN, key = getServerFlags()

	_ = env.Parse(&cfg)

	flag.Parse()

	if cfg.Host != "" {
		host = &cfg.Host
	}

	//STORE_INTERVAL can be set to 0, hence can't check it as !=0
	value, ok := os.LookupEnv("STORE_INTERVAL")
	if ok && value != "" {
		storeIntervalIn = &cfg.StoreInterval
	}

	if cfg.FileStoragePath != "" {
		fileStoragePath = &cfg.FileStoragePath
	}

	if cfg.DatabaseDSN != "" {
		dbDSN = &cfg.DatabaseDSN
	}

	if cfg.Key != "" {
		key = &cfg.Key
	}

	value, ok = os.LookupEnv("RESTORE")
	if ok && value != "" {
		restore = &cfg.Restore
	}

	storeIntervalOut = utils.GetTimeInterval(*storeIntervalIn)

	return
}

// host=localhost user=postgres password=postgres sslmode=disable dbname=evgenytrefilov
func getServerFlags() (host *string, storeInterval *float64, fileStoragePath *string, restore *bool, dbDSN, key *string) {
	host = flag.String("a", "localhost:8080", "host address")
	storeInterval = flag.Float64("i", 300, "file store interval")
	fileStoragePath = flag.String("f", "/tmp/metrics-db.json", "file storage path")
	dbDSN = flag.String("d", "", "database address")
	restore = flag.Bool("r", true, "restore saved metrics on server start")
	key = flag.String("k", "", "hash key")
	return
}

func getAgentFlags() (host *string, pollInterval, reportInterval *float64, key *string, rateLimit *int64) {
	host = flag.String("a", "localhost:8080", "host address")
	pollInterval = flag.Float64("p", 2, "metrics polling interval")
	reportInterval = flag.Float64("r", 10, "metrics reporting interval")
	key = flag.String("k", "", "hash key")
	rateLimit = flag.Int64("l", 2, "metrics report rate limit")
	return
}
