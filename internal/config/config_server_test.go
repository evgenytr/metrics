package config

import (
	"testing"
	"time"
)

func TestGetServerConfig(t *testing.T) {
	tests := []struct {
		name                 string
		wantHost             string
		wantHostGrpc         string
		wantStoreIntervalOut time.Duration
		wantFileStoragePath  string
		wantRestore          bool
		wantDBDSN            string
		wantKey              string
		wantCryptoKey        string
		wantTrustedSubnet    string
	}{
		{
			name:                 "Defaults",
			wantHost:             "localhost:8080",
			wantHostGrpc:         "localhost:3200",
			wantStoreIntervalOut: 300 * time.Second,
			wantFileStoragePath:  "/tmp/metrics-db.json",
			wantRestore:          true,
			wantDBDSN:            "",
			wantKey:              "",
			wantCryptoKey:        "./rsakeys/private.pem",
			wantTrustedSubnet:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, gotHostGrpc, gotStoreIntervalOut, gotFileStoragePath, gotRestore, gotDBDSN, gotKey, gotCryptoKey, gotTrustedSubnet := GetServerConfig()
			if gotHost != tt.wantHost {
				t.Errorf("GetServerConfig() gotHost = %v, want %v", gotHost, tt.wantHost)
			}
			if gotHostGrpc != tt.wantHostGrpc {
				t.Errorf("GetServerConfig() gotHostGrpc = %v, want %v", gotHostGrpc, tt.wantHostGrpc)
			}
			if gotStoreIntervalOut != tt.wantStoreIntervalOut {
				t.Errorf("GetServerConfig() gotStoreIntervalOut = %v, want %v", gotStoreIntervalOut, tt.wantStoreIntervalOut)
			}
			if gotFileStoragePath != tt.wantFileStoragePath {
				t.Errorf("GetServerConfig() gotFileStoragePath = %v, want %v", gotFileStoragePath, tt.wantFileStoragePath)
			}
			if gotRestore != tt.wantRestore {
				t.Errorf("GetServerConfig() gotRestore = %v, want %v", gotRestore, tt.wantRestore)
			}
			if gotDBDSN != tt.wantDBDSN {
				t.Errorf("GetServerConfig() gotDBDSN = %v, want %v", gotDBDSN, tt.wantDBDSN)
			}
			if gotKey != tt.wantKey {
				t.Errorf("GetServerConfig() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
			if gotCryptoKey != tt.wantCryptoKey {
				t.Errorf("GetServerConfig() gotCryptoKey = %v, want %v", gotCryptoKey, tt.wantCryptoKey)
			}
			if gotTrustedSubnet != tt.wantTrustedSubnet {
				t.Errorf("GetServerConfig() gotTrustedSubnet = %v, want %v", gotTrustedSubnet, tt.wantTrustedSubnet)
			}
		})
	}
}
