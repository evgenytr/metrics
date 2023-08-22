package logging

import (
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestNewLogger(t *testing.T) {
	type args struct {
		config string
	}
	tests := []struct {
		name       string
		args       args
		wantLogger *zap.Logger
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLogger, err := NewLogger(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotLogger, tt.wantLogger) {
				t.Errorf("NewLogger() gotLogger = %v, want %v", gotLogger, tt.wantLogger)
			}
		})
	}
}
