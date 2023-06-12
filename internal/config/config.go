package config

import (
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

type agentConfig struct {
	Host           string  `env:"ADDRESS"`
	ReportInterval float64 `env:"REPORT_INTERVAL"`
	PollInterval   float64 `env:"POLL_INTERVAL"`
}

type serverConfig struct {
	Host            string  `env:"ADDRESS"`
	StoreInterval   float64 `env:"STORE_INTERVAL"`
	FileStoragePath string  `env:"FILE_STORAGE_PATH"`
	Restore         bool    `env:"RESTORE"`
	DatabaseDSN     string  `env:"DATABASE_DSN"`
}

func GetAgentConfig() (host *string, pollInterval, reportInterval *float64) {
	host, pollInterval, reportInterval = getAgentFlags()
	var cfg agentConfig
	_ = env.Parse(&cfg)
	flag.Parse()

	if cfg.Host != "" {
		host = &cfg.Host
	}

	if cfg.PollInterval != 0 {
		pollInterval = &cfg.PollInterval
	}

	if cfg.ReportInterval != 0 {
		reportInterval = &cfg.ReportInterval
	}
	return
}
func GetServerConfig() (host *string, storeIntervalOut *time.Duration, fileStoragePath *string, restore *bool, dbDSN *string) {

	var storeIntervalIn *float64
	var cfg serverConfig

	host, storeIntervalIn, fileStoragePath, restore, dbDSN = getServerFlags()

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

	value, ok = os.LookupEnv("RESTORE")
	if ok && value != "" {
		restore = &cfg.Restore
	}

	storeIntervalValue := time.Duration(*storeIntervalIn) * time.Second
	storeIntervalOut = &storeIntervalValue
	return
}

// host=localhost user=postgres password=postgres sslmode=disable
func getServerFlags() (host *string, storeInterval *float64, fileStoragePath *string, restore *bool, dbDSN *string) {
	host = flag.String("a", "localhost:8080", "host address")
	storeInterval = flag.Float64("i", 300, "file store interval")
	fileStoragePath = flag.String("f", "/tmp/metrics-db.json", "file storage path")
	dbDSN = flag.String("d", "", "database address")
	restore = flag.Bool("r", true, "restore saved metrics on server start")
	return
}

func getAgentFlags() (host *string, pollInterval, reportInterval *float64) {
	host = flag.String("a", "localhost:8080", "host address")
	pollInterval = flag.Float64("p", 2, "metrics polling interval")
	reportInterval = flag.Float64("r", 10, "metrics reporting interval")
	return
}
