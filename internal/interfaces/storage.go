// Package interfaces contains interfaces used across the services
package interfaces

import (
	"context"

	"github.com/evgenytr/metrics.git/internal/metric"
)

// Storage interface should be implemented by any used storage
type Storage interface {
	InitializeMetrics(ctx context.Context, restore bool) error
	StoreMetrics(ctx context.Context) error
	UpdateGauge(ctx context.Context, name string, value float64) (float64, error)
	UpdateCounter(ctx context.Context, name string, value int64) (int64, error)
	GetGaugeValue(ctx context.Context, name string) (float64, error)
	GetCounterValue(ctx context.Context, name string) (int64, error)
	Update(ctx context.Context, metricType, name, value string) (string, error)
	ReadValue(ctx context.Context, metricType, name string) (string, error)
	ListAll(ctx context.Context) (map[string]string, error)
	GetMetricsMap(ctx context.Context) (map[string]*metric.Metrics, error)
	Ping(ctx context.Context) error
}
