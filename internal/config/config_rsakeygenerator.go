// Package config populates RSA key generator config data
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type rsaKeyGeneratorConfig struct {
	PublicKey  string `json:"public_key_path,omitempty"`
	PrivateKey string `json:"private_key_path,omitempty"`
}

// GetRsaKeyGeneratorConfig returns RSA key generator config params
func GetRsaKeyGeneratorConfig() (publicKeyPath, privateKeyPath string, err error) {

	publicKeyPath, privateKeyPath, err = getRsaKeyGeneratorFlags()

	return
}

func getRsaKeyGeneratorFlags() (publicKeyPath, privateKeyPath string, err error) {

	var configFile string
	flag.StringVar(&configFile, "config", "", "config JSON file path")

	flag.Parse()

	if configFile == "" {
		err = fmt.Errorf("configFile flag is missing or empty")
		return
	}

	configDefaults := &rsaKeyGeneratorConfig{}

	fmt.Println("RSA key generator config file ", configFile)

	fmt.Println("reading config file")
	dat, err := os.ReadFile(configFile)
	if err != nil {
		err = fmt.Errorf("failed to read config file %w", err)
		return
	}
	fmt.Println(dat)
	err = json.Unmarshal(dat, configDefaults)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal config file %w", err)
		return
	}
	fmt.Println(configDefaults)

	publicKeyPath = configDefaults.PublicKey
	privateKeyPath = configDefaults.PrivateKey

	return
}
