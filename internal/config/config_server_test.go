package config

import (
	"reflect"
	"testing"
	"time"
)

func TestGetServerConfig(t *testing.T) {
	tests := []struct {
		name                 string
		wantHost             string
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
			wantStoreIntervalOut: 300 * time.Second,
			wantFileStoragePath:  "/tmp/metrics-db.json",
			wantRestore:          true,
			wantDBDSN:            "",
			wantKey:              "",
			wantCryptoKey:        "./rsakeys/private.pem",
			wantTrustedSubnet:    "192.168.0.0/24",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, gotStoreIntervalOut, gotFileStoragePath, gotRestore, gotDBDSN, gotKey, gotCryptoKey, gotTrustedSubnet := GetServerConfig()
			if !reflect.DeepEqual(gotHost, tt.wantHost) {
				t.Errorf("GetServerConfig() gotHost = %v, want %v", gotHost, tt.wantHost)
			}
			if !reflect.DeepEqual(gotStoreIntervalOut, tt.wantStoreIntervalOut) {
				t.Errorf("GetServerConfig() gotStoreIntervalOut = %v, want %v", gotStoreIntervalOut, tt.wantStoreIntervalOut)
			}
			if !reflect.DeepEqual(gotFileStoragePath, tt.wantFileStoragePath) {
				t.Errorf("GetServerConfig() gotFileStoragePath = %v, want %v", gotFileStoragePath, tt.wantFileStoragePath)
			}
			if !reflect.DeepEqual(gotRestore, tt.wantRestore) {
				t.Errorf("GetServerConfig() gotRestore = %v, want %v", gotRestore, tt.wantRestore)
			}
			if !reflect.DeepEqual(gotDBDSN, tt.wantDBDSN) {
				t.Errorf("GetServerConfig() gotDBDSN = %v, want %v", gotDBDSN, tt.wantDBDSN)
			}
			if !reflect.DeepEqual(gotKey, tt.wantKey) {
				t.Errorf("GetServerConfig() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
			if !reflect.DeepEqual(gotCryptoKey, tt.wantCryptoKey) {
				t.Errorf("GetAgentConfig() gotCryptoKey = %v, want %v", gotCryptoKey, tt.wantCryptoKey)
			}
			if !reflect.DeepEqual(gotTrustedSubnet, tt.wantTrustedSubnet) {
				t.Errorf("GetAgentConfig() gotCryptoKey = %v, want %v", gotCryptoKey, tt.wantCryptoKey)
			}
		})
	}
}
