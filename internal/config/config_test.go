package config

import (
	"reflect"
	"testing"
	"time"
)

func TestGetAgentConfig(t *testing.T) {
	tests := []struct {
		name                  string
		wantHost              *string
		wantPollIntervalOut   *time.Duration
		wantReportIntervalOut *time.Duration
		wantKey               *string
		wantRateLimit         *int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, gotPollIntervalOut, gotReportIntervalOut, gotKey, gotRateLimit := GetAgentConfig()
			if !reflect.DeepEqual(gotHost, tt.wantHost) {
				t.Errorf("GetAgentConfig() gotHost = %v, want %v", gotHost, tt.wantHost)
			}
			if !reflect.DeepEqual(gotPollIntervalOut, tt.wantPollIntervalOut) {
				t.Errorf("GetAgentConfig() gotPollIntervalOut = %v, want %v", gotPollIntervalOut, tt.wantPollIntervalOut)
			}
			if !reflect.DeepEqual(gotReportIntervalOut, tt.wantReportIntervalOut) {
				t.Errorf("GetAgentConfig() gotReportIntervalOut = %v, want %v", gotReportIntervalOut, tt.wantReportIntervalOut)
			}
			if !reflect.DeepEqual(gotKey, tt.wantKey) {
				t.Errorf("GetAgentConfig() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
			if !reflect.DeepEqual(gotRateLimit, tt.wantRateLimit) {
				t.Errorf("GetAgentConfig() gotRateLimit = %v, want %v", gotRateLimit, tt.wantRateLimit)
			}
		})
	}
}

func TestGetServerConfig(t *testing.T) {
	tests := []struct {
		name                 string
		wantHost             *string
		wantStoreIntervalOut *time.Duration
		wantFileStoragePath  *string
		wantRestore          *bool
		wantDBDSN            *string
		wantKey              *string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, gotStoreIntervalOut, gotFileStoragePath, gotRestore, gotDBDSN, gotKey := GetServerConfig()
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
		})
	}
}

func Test_getAgentFlags(t *testing.T) {
	tests := []struct {
		name               string
		wantHost           *string
		wantPollInterval   *float64
		wantReportInterval *float64
		wantKey            *string
		wantRateLimit      *int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, gotPollInterval, gotReportInterval, gotKey, gotRateLimit := getAgentFlags()
			if !reflect.DeepEqual(gotHost, tt.wantHost) {
				t.Errorf("getAgentFlags() gotHost = %v, want %v", gotHost, tt.wantHost)
			}
			if !reflect.DeepEqual(gotPollInterval, tt.wantPollInterval) {
				t.Errorf("getAgentFlags() gotPollInterval = %v, want %v", gotPollInterval, tt.wantPollInterval)
			}
			if !reflect.DeepEqual(gotReportInterval, tt.wantReportInterval) {
				t.Errorf("getAgentFlags() gotReportInterval = %v, want %v", gotReportInterval, tt.wantReportInterval)
			}
			if !reflect.DeepEqual(gotKey, tt.wantKey) {
				t.Errorf("getAgentFlags() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
			if !reflect.DeepEqual(gotRateLimit, tt.wantRateLimit) {
				t.Errorf("getAgentFlags() gotRateLimit = %v, want %v", gotRateLimit, tt.wantRateLimit)
			}
		})
	}
}

func Test_getServerFlags(t *testing.T) {
	tests := []struct {
		name                string
		wantHost            *string
		wantStoreInterval   *float64
		wantFileStoragePath *string
		wantRestore         *bool
		wantDBDSN           *string
		wantKey             *string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, gotStoreInterval, gotFileStoragePath, gotRestore, gotDBDSN, gotKey := getServerFlags()
			if !reflect.DeepEqual(gotHost, tt.wantHost) {
				t.Errorf("getServerFlags() gotHost = %v, want %v", gotHost, tt.wantHost)
			}
			if !reflect.DeepEqual(gotStoreInterval, tt.wantStoreInterval) {
				t.Errorf("getServerFlags() gotStoreInterval = %v, want %v", gotStoreInterval, tt.wantStoreInterval)
			}
			if !reflect.DeepEqual(gotFileStoragePath, tt.wantFileStoragePath) {
				t.Errorf("getServerFlags() gotFileStoragePath = %v, want %v", gotFileStoragePath, tt.wantFileStoragePath)
			}
			if !reflect.DeepEqual(gotRestore, tt.wantRestore) {
				t.Errorf("getServerFlags() gotRestore = %v, want %v", gotRestore, tt.wantRestore)
			}
			if !reflect.DeepEqual(gotDBDSN, tt.wantDBDSN) {
				t.Errorf("getServerFlags() gotDBDSN = %v, want %v", gotDBDSN, tt.wantDBDSN)
			}
			if !reflect.DeepEqual(gotKey, tt.wantKey) {
				t.Errorf("getServerFlags() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}
