// Package config populates metric agent and server config data based on flags, environment variables or defaults
package config

import (
	"encoding/json"
	"flag"
	"fmt"
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
	TrustedSubnet   string  `env:"TRUSTED_SUBNET" json:"trusted_subnet,omitempty"`
	ConfigFile      string  `env:"CONFIG"`
}

// GetServerConfig returns server config params
func GetServerConfig() (host string, storeIntervalOut time.Duration, fileStoragePath string, restore bool, dbDSN, key, cryptoKey, trustedSubnet string) {

	var storeIntervalIn float64
	var cfg serverConfig

	_ = env.Parse(&cfg)

	host, storeIntervalIn, fileStoragePath, restore, dbDSN, key, cryptoKey, trustedSubnet = getServerFlags(cfg.ConfigFile)

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

	if cfg.TrustedSubnet != "" {
		trustedSubnet = cfg.TrustedSubnet
	}

	storeIntervalOut = utils.GetTimeInterval(storeIntervalIn)

	return
}

func getServerFlags(configFile string) (host string, storeInterval float64, fileStoragePath string, restore bool, dbDSN, key, cryptoKey, trustedSubnet string) {

	if flag.Lookup("config") == nil && configFile == "" {
		fmt.Println("setting config file from flag")
		flag.StringVar(&configFile, "config", "", "config JSON file path")
	}

	if flag.Lookup("a") == nil {
		flag.StringVar(&host, "a", "", "host address")
	}
	if flag.Lookup("i") == nil {
		flag.Float64Var(&storeInterval, "i", -1, "file store interval")
	}
	if flag.Lookup("f") == nil {
		flag.StringVar(&fileStoragePath, "f", "", "file storage path")
	}
	if flag.Lookup("d") == nil {
		flag.StringVar(&dbDSN, "d", "", "database address")
	}
	if flag.Lookup("r") == nil {
		flag.BoolVar(&restore, "r", false, "restore saved metrics on server start")
	}
	if flag.Lookup("k") == nil {
		flag.StringVar(&key, "k", "", "hash key")
	}
	if flag.Lookup("crypto-key") == nil {
		flag.StringVar(&cryptoKey, "crypto-key", "", "crypto key path")
	}
	if flag.Lookup("t") == nil {
		flag.StringVar(&trustedSubnet, "t", "", "trusted subnet CIDR")
	}
	flag.Parse()

	//sensible defaults to run in absence of flags and env vars
	configDefaults := &serverConfig{
		Host:            "localhost:8080",
		StoreInterval:   300,
		FileStoragePath: "/tmp/metrics-db.json",
		Restore:         true,
		CryptoKey:       "./rsakeys/private.pem",
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

	if key == "" {
		key = configDefaults.Key
	}

	if cryptoKey == "" {
		cryptoKey = configDefaults.CryptoKey
	}

	if storeInterval == -1 {
		storeInterval = configDefaults.StoreInterval
	}

	if !restore {
		restore = configDefaults.Restore
	}

	if dbDSN == "" {
		dbDSN = configDefaults.DatabaseDSN
	}

	if fileStoragePath == "" {
		fileStoragePath = configDefaults.FileStoragePath
	}

	if trustedSubnet == "" {
		trustedSubnet = configDefaults.TrustedSubnet
	}

	return
}
