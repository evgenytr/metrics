package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"os"
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
func GetServerConfig() (host *string, storeInterval *float64, fileStoragePath *string, restore *bool) {

	host, storeInterval, fileStoragePath, restore = getServerFlags()
	var cfg serverConfig

	_ = env.Parse(&cfg)

	flag.Parse()

	if cfg.Host != "" {
		host = &cfg.Host
	}

	//STORE_INTERVAL can be set to 0, hence can't check it as !=0
	value, ok := os.LookupEnv("STORE_INTERVAL")
	if ok && value != "" {
		storeInterval = &cfg.StoreInterval
	}

	if cfg.FileStoragePath != "" {
		fileStoragePath = &cfg.FileStoragePath
	}

	value, ok = os.LookupEnv("RESTORE")
	if ok && value != "" {
		restore = &cfg.Restore
	}

	return
}

func getServerFlags() (host *string, storeInterval *float64, fileStoragePath *string, restore *bool) {
	host = flag.String("a", "localhost:8080", "host address")
	storeInterval = flag.Float64("i", 300, "file store interval")
	fileStoragePath = flag.String("f", "/tmp/metrics-db.json", "file storage path")
	restore = flag.Bool("r", true, "restore saved metrics on server start")
	return
}

func getAgentFlags() (host *string, pollInterval, reportInterval *float64) {
	host = flag.String("a", "localhost:8080", "host address")
	pollInterval = flag.Float64("p", 2, "metrics polling interval")
	reportInterval = flag.Float64("r", 10, "metrics reporting interval")
	return
}
