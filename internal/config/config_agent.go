// Package config populates metric agent and server config data based on flags, environment variables or defaults
package config

import (
	"flag"
	"time"

	"github.com/evgenytr/metrics.git/internal/utils"

	"github.com/caarlos0/env/v6"
)

type agentConfig struct {
	Host           string  `env:"ADDRESS" json:"address,omitempty"`
	Key            string  `env:"KEY" json:"key,omitempty"`
	ReportInterval float64 `env:"REPORT_INTERVAL" json:"report_interval,omitempty"`
	PollInterval   float64 `env:"POLL_INTERVAL" json:"poll_interval,omitempty"`
	RateLimit      int64   `env:"RATE_LIMIT" json:"rate_limit,omitempty"`
	CryptoKey      string  `env:"CRYPTO_KEY" json:"crypto_key,omitempty"`
	ConfigFile     string  `env:"CONFIG"`
}

// GetAgentConfig returns agent config params
func GetAgentConfig() (host string, pollIntervalOut, reportIntervalOut time.Duration, key string, rateLimit int64, cryptoKeyFile string) {

	var cfg agentConfig
	var pollIntervalIn, reportIntervalIn float64

	_ = env.Parse(&cfg)
	flag.Parse()

	host, pollIntervalIn, reportIntervalIn, key, rateLimit, cryptoKeyFile = getAgentFlags(cfg.ConfigFile)

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

	if cfg.CryptoKey != "" {
		cryptoKeyFile = cfg.CryptoKey
	}

	pollIntervalOut = utils.GetTimeInterval(pollIntervalIn)
	reportIntervalOut = utils.GetTimeInterval(reportIntervalIn)

	return
}

func getAgentFlags(configFile string) (host string, pollInterval, reportInterval float64, key string, rateLimit int64, cryptoKeyFile string) {

	configDefaults := &agentConfig{
		Host:           "localhost:8080",
		PollInterval:   2,
		ReportInterval: 10,
		RateLimit:      2,
		CryptoKey:      "./rsakeys/public.pub",
	}
	if flag.Lookup("config") == nil && configFile == "" {
		configFile = *flag.String("config", "", "config JSON file path")
	}

	//TODO: if config file exists, set defaults based on it

	if flag.Lookup("a") == nil {
		host = *flag.String("a", configDefaults.Host, "host address")
	}
	if flag.Lookup("p") == nil {
		pollInterval = *flag.Float64("p", configDefaults.PollInterval, "metrics polling interval")
	}
	if flag.Lookup("r") == nil {
		reportInterval = *flag.Float64("r", configDefaults.ReportInterval, "metrics reporting interval")
	}
	if flag.Lookup("k") == nil {
		key = *flag.String("k", configDefaults.Key, "hash key")
	}
	if flag.Lookup("l") == nil {
		rateLimit = *flag.Int64("l", configDefaults.RateLimit, "metrics report rate limit")
	}
	if flag.Lookup("crypto-key") == nil {
		cryptoKeyFile = *flag.String("crypto-key", configDefaults.CryptoKey, "crypto key file path")
	}

	return
}
