package config

import (
	"reflect"
	"testing"
	"time"
)

func TestGetAgentConfig(t *testing.T) {
	tests := []struct {
		name                  string
		wantHost              string
		wantPollIntervalOut   time.Duration
		wantReportIntervalOut time.Duration
		wantKey               string
		wantRateLimit         int64
	}{
		{
			name:                  "Defaults",
			wantHost:              "localhost:8080",
			wantPollIntervalOut:   2 * time.Second,
			wantReportIntervalOut: 10 * time.Second,
			wantKey:               "",
			wantRateLimit:         2,
		},
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
