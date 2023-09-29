// Package config populates metric agent and server config data based on flags, environment variables or defaults
package config

import (
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env/v6"

	"github.com/evgenytr/metrics.git/internal/utils"
)

type serverConfig struct {
	Host            string  `env:"ADDRESS" json:"address,omitempty"`
	FileStoragePath string  `env:"FILE_STORAGE_PATH" json:"store_file,omitempty"`
	DatabaseDSN     string  `env:"DATABASE_DSN" json:"database_dsn,omitempty"`
	Key             string  `env:"KEY" json:"key,omitempty"`
	StoreInterval   float64 `env:"STORE_INTERVAL" json:"store_interval,omitempty"`
	Restore         bool    `env:"RESTORE" json:"restore,omitempty"`
	CryptoKey       string  `env:"CRYPTO_KEY" json:"crypto_key,omitempty"`
	ConfigFile      string  `env:"CONFIG"`
}

// GetServerConfig returns server config params
func GetServerConfig() (host string, storeIntervalOut time.Duration, fileStoragePath string, restore bool, dbDSN, key, cryptoKey string) {

	var storeIntervalIn float64
	var cfg serverConfig

	_ = env.Parse(&cfg)
	flag.Parse()

	host, storeIntervalIn, fileStoragePath, restore, dbDSN, key, cryptoKey = getServerFlags(cfg.ConfigFile)

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

func getServerFlags(configFile string) (host string, storeInterval float64, fileStoragePath string, restore bool, dbDSN, key, cryptoKey string) {
	host = "localhost:8080"
	restore = true

	configDefaults := &serverConfig{
		Host:            "localhost:8080",
		StoreInterval:   300,
		FileStoragePath: "/tmp/metrics-db.json",
		Restore:         true,
		CryptoKey:       "./rsakeys/private.pem",
	}
	if flag.Lookup("config") == nil && configFile == "" {
		configFile = *flag.String("config", "", "config JSON file path")
	}

	//TODO: if config file exists, set defaults based on it

	if flag.Lookup("a") == nil {
		host = *flag.String("a", configDefaults.Host, "host address")
	}
	if flag.Lookup("i") == nil {
		storeInterval = *flag.Float64("i", configDefaults.StoreInterval, "file store interval")
	}
	if flag.Lookup("f") == nil {
		fileStoragePath = *flag.String("f", configDefaults.FileStoragePath, "file storage path")
	}
	if flag.Lookup("d") == nil {
		dbDSN = *flag.String("d", "", "database address")
	}
	if flag.Lookup("r") == nil {
		restore = *flag.Bool("r", configDefaults.Restore, "restore saved metrics on server start")
	}
	if flag.Lookup("k") == nil {
		key = *flag.String("k", "", "hash key")
	}
	if flag.Lookup("crypto-key") == nil {
		cryptoKey = *flag.String("crypto-key", configDefaults.CryptoKey, "crypto key path")
	}
	return
}
