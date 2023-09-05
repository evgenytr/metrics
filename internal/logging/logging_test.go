package logging

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	type args struct {
		config string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Development",
			args: args{
				config: NewDevelopment,
			},
			wantErr: false,
		},
		{
			name: "Production",
			args: args{
				config: NewProduction,
			},
			wantErr: false,
		},
		{
			name: "Default",
			args: args{
				config: "Any string returns Development logger",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			_, err := NewLogger(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
