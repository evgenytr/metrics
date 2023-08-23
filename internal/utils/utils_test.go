package utils

import (
	"math/rand"
	"testing"
)

func BenchmarkGetTimeInterval(b *testing.B) {
	for i := 0; i < b.N; i++ {
		interval := rand.Float64() * 100
		GetTimeInterval(interval)
	}
}

/*
func TestGetTimeInterval(t *testing.T) {
	type args struct {
		seconds float64
	}
	tests := []struct {
		name string
		args args
		want *time.Duration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTimeInterval(tt.args.seconds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTimeInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
