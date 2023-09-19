// Package config populates metric agent and server config data based on flags, environment variables or defaults
package config

import (
	"flag"
	"os"
	"time"

	"github.com/evgenytr/metrics.git/internal/utils"

	"github.com/caarlos0/env/v6"
)

type serverConfig struct {
	Host            string  `env:"ADDRESS"`
	FileStoragePath string  `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string  `env:"DATABASE_DSN"`
	Key             string  `env:"KEY"`
	StoreInterval   float64 `env:"STORE_INTERVAL"`
	Restore         bool    `env:"RESTORE"`
	CryptoKey       string  `env:"CRYPTO_KEY"`
}

// GetServerConfig returns server config params
func GetServerConfig() (host string, storeIntervalOut time.Duration, fileStoragePath string, restore bool, dbDSN, key, cryptoKey string) {

	var storeIntervalIn float64
	var cfg serverConfig

	host, storeIntervalIn, fileStoragePath, restore, dbDSN, key, cryptoKey = getServerFlags()

	_ = env.Parse(&cfg)

	flag.Parse()

	if cfg.Host != "" {
		host = cfg.Host
	}

	//STORE_INTERVAL can be set to 0, hence can't check it as !=0
	value, ok := os.LookupEnv("STORE_INTERVAL")
	if ok && value != "" {
		storeIntervalIn = cfg.StoreInterval
	}

	if cfg.FileStoragePath != "" {
		fileStoragePath = cfg.FileStoragePath
	}

	if cfg.DatabaseDSN != "" {
		dbDSN = cfg.DatabaseDSN
	}

	if cfg.Key != "" {
		key = cfg.Key
	}

	value, ok = os.LookupEnv("RESTORE")
	if ok && value != "" {
		restore = cfg.Restore
	}

	if cfg.CryptoKey != "" {
		cryptoKey = cfg.CryptoKey
	}

	storeIntervalOut = utils.GetTimeInterval(storeIntervalIn)

	return
}

func getServerFlags() (host string, storeInterval float64, fileStoragePath string, restore bool, dbDSN, key, cryptoKey string) {
	host = "localhost:8080"
	restore = true
	if flag.Lookup("a") == nil {
		host = *flag.String("a", "localhost:8080", "host address")
	}
	if flag.Lookup("i") == nil {
		storeInterval = *flag.Float64("i", 300, "file store interval")
	}
	if flag.Lookup("f") == nil {
		fileStoragePath = *flag.String("f", "/tmp/metrics-db.json", "file storage path")
	}
	if flag.Lookup("d") == nil {
		dbDSN = *flag.String("d", "", "database address")
	}
	if flag.Lookup("r") == nil {
		restore = *flag.Bool("r", true, "restore saved metrics on server start")
	}
	if flag.Lookup("k") == nil {
		key = *flag.String("k", "", "hash key")
	}
	if flag.Lookup("crypto-key") == nil {
		cryptoKey = *flag.String("crypto-key", "", "crypto key path")
	}
	return
}
