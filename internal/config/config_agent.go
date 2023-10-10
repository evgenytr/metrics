// Package config populates metric agent and server config data based on flags, environment variables or defaults
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/evgenytr/metrics.git/internal/utils"

	"github.com/caarlos0/env/v6"
)

type agentConfig struct {
	Host           string  `env:"ADDRESS" json:"address,omitempty"`
	HostGrpc       string  `json:"address_grpc,omitempty"`
	Key            string  `env:"KEY" json:"key,omitempty"`
	ReportInterval float64 `env:"REPORT_INTERVAL" json:"report_interval,omitempty"`
	PollInterval   float64 `env:"POLL_INTERVAL" json:"poll_interval,omitempty"`
	RateLimit      int64   `env:"RATE_LIMIT" json:"rate_limit,omitempty"`
	CryptoKey      string  `env:"CRYPTO_KEY" json:"crypto_key,omitempty"`
	ConfigFile     string  `env:"CONFIG"`
}

// GetAgentConfig returns agent config params
func GetAgentConfig() (host, hostGrpc string, pollIntervalOut, reportIntervalOut time.Duration, key string, rateLimit int64, cryptoKeyFile string) {

	var cfg agentConfig
	var pollIntervalIn, reportIntervalIn float64

	_ = env.Parse(&cfg)

	fmt.Println(cfg)

	host, hostGrpc, pollIntervalIn, reportIntervalIn, key, rateLimit, cryptoKeyFile = getAgentFlags(cfg.ConfigFile)

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

func getAgentFlags(configFile string) (host, hostGrpc string, pollInterval, reportInterval float64, key string, rateLimit int64, cryptoKeyFile string) {

	if flag.Lookup("config") == nil && configFile == "" {
		fmt.Println("setting config file from flag")
		flag.StringVar(&configFile, "config", "", "config JSON file path")
	}

	if flag.Lookup("a") == nil {
		flag.StringVar(&host, "a", "", "host address")
	}
	if flag.Lookup("p") == nil {
		flag.Float64Var(&pollInterval, "p", -1, "metrics polling interval")
	}
	if flag.Lookup("r") == nil {
		flag.Float64Var(&reportInterval, "r", -1, "metrics reporting interval")
	}
	if flag.Lookup("k") == nil {
		flag.StringVar(&key, "k", "", "hash key")
	}
	if flag.Lookup("l") == nil {
		flag.Int64Var(&rateLimit, "l", -1, "metrics report rate limit")
	}
	if flag.Lookup("crypto-key") == nil {
		flag.StringVar(&cryptoKeyFile, "crypto-key", "", "crypto key file path")
	}

	flag.Parse()

	//sensible defaults to run in absence of flags and env vars
	configDefaults := &agentConfig{
		Host:           "localhost:8080",
		HostGrpc:       "localhost:3200",
		PollInterval:   1,
		ReportInterval: 10,
		RateLimit:      2,
		CryptoKey:      "./rsakeys/public.pub",
	}

	fmt.Println("config file ", configFile)

	if configFile != "" {
		fmt.Println("reading config file")
		dat, err := os.ReadFile(configFile)
		if err != nil {
			fmt.Println("failed to read config file %w", err)
		}
		fmt.Println(dat)
		err = json.Unmarshal(dat, configDefaults)
		if err != nil {
			fmt.Println("failed to unmarshal config file %w", err)
		}
		fmt.Println(configDefaults)

	}

	//set defaults for values that were not set from flags
	if host == "" {
		host = configDefaults.Host
	}

	if hostGrpc == "" {
		hostGrpc = configDefaults.HostGrpc
	}

	if key == "" {
		key = configDefaults.Key
	}

	if cryptoKeyFile == "" {
		cryptoKeyFile = configDefaults.CryptoKey
	}

	if rateLimit == -1 {
		rateLimit = configDefaults.RateLimit
	}

	if pollInterval == -1 {
		pollInterval = configDefaults.PollInterval
	}

	if reportInterval == -1 {
		reportInterval = configDefaults.ReportInterval
	}

	return
}
