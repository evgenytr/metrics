package utils

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func BenchmarkGetTimeInterval(b *testing.B) {
	for i := 0; i < b.N; i++ {
		interval := rand.Float64() * 100
		GetTimeInterval(interval)
	}
}

func TestGetTimeInterval(t *testing.T) {
	type args struct {
		seconds float64
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "10 seconds",
			args: args{
				seconds: 10,
			},
			want: time.Second * 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTimeInterval(tt.args.seconds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTimeInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}
