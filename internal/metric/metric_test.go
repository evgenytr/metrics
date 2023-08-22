package metric

import (
	"reflect"
	"testing"
)

func TestCreate(t *testing.T) {
	type args struct {
		metricType string
		name       string
		value      string
	}
	tests := []struct {
		name          string
		args          args
		wantNewMetric *Metrics
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewMetric, err := Create(tt.args.metricType, tt.args.name, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNewMetric, tt.wantNewMetric) {
				t.Errorf("Create() gotNewMetric = %v, want %v", gotNewMetric, tt.wantNewMetric)
			}
		})
	}
}

func TestCreateCounter(t *testing.T) {
	type args struct {
		name  string
		value *int64
	}
	tests := []struct {
		name          string
		args          args
		wantNewMetric *Metrics
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewMetric, err := CreateCounter(tt.args.name, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNewMetric, tt.wantNewMetric) {
				t.Errorf("CreateCounter() gotNewMetric = %v, want %v", gotNewMetric, tt.wantNewMetric)
			}
		})
	}
}

func TestCreateGauge(t *testing.T) {
	type args struct {
		name  string
		value *float64
	}
	tests := []struct {
		name          string
		args          args
		wantNewMetric *Metrics
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewMetric, err := CreateGauge(tt.args.name, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateGauge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNewMetric, tt.wantNewMetric) {
				t.Errorf("CreateGauge() gotNewMetric = %v, want %v", gotNewMetric, tt.wantNewMetric)
			}
		})
	}
}

func TestMetrics_Add(t *testing.T) {
	type fields struct {
		ID    string
		MType string
		Delta *int64
		Value *float64
	}
	type args struct {
		metricType string
		value      string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantNewValue string
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric := &Metrics{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
			}
			gotNewValue, err := metric.Add(tt.args.metricType, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNewValue != tt.wantNewValue {
				t.Errorf("Add() gotNewValue = %v, want %v", gotNewValue, tt.wantNewValue)
			}
		})
	}
}

func TestMetrics_GetCounterValue(t *testing.T) {
	type fields struct {
		ID    string
		MType string
		Delta *int64
		Value *float64
	}
	tests := []struct {
		name      string
		fields    fields
		wantValue *int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric := &Metrics{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
			}
			if gotValue := metric.GetCounterValue(); !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("GetCounterValue() = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

func TestMetrics_GetGaugeValue(t *testing.T) {
	type fields struct {
		ID    string
		MType string
		Delta *int64
		Value *float64
	}
	tests := []struct {
		name      string
		fields    fields
		wantValue *float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric := &Metrics{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
			}
			if gotValue := metric.GetGaugeValue(); !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("GetGaugeValue() = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

func TestMetrics_GetType(t *testing.T) {
	type fields struct {
		ID    string
		MType string
		Delta *int64
		Value *float64
	}
	tests := []struct {
		name      string
		fields    fields
		wantValue string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric := &Metrics{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
			}
			if gotValue := metric.GetType(); gotValue != tt.wantValue {
				t.Errorf("GetType() = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

func TestMetrics_GetValue(t *testing.T) {
	type fields struct {
		ID    string
		MType string
		Delta *int64
		Value *float64
	}
	tests := []struct {
		name      string
		fields    fields
		wantValue string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric := &Metrics{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
			}
			if gotValue := metric.GetValue(); gotValue != tt.wantValue {
				t.Errorf("GetValue() = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

func TestMetrics_ResetCounter(t *testing.T) {
	type fields struct {
		ID    string
		MType string
		Delta *int64
		Value *float64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric := &Metrics{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
			}
			if err := metric.ResetCounter(); (err != nil) != tt.wantErr {
				t.Errorf("ResetCounter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMetrics_UpdateCounter(t *testing.T) {
	type fields struct {
		ID    string
		MType string
		Delta *int64
		Value *float64
	}
	type args struct {
		value *int64
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantNewValue *int64
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric := &Metrics{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
			}
			gotNewValue, err := metric.UpdateCounter(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNewValue, tt.wantNewValue) {
				t.Errorf("UpdateCounter() gotNewValue = %v, want %v", gotNewValue, tt.wantNewValue)
			}
		})
	}
}

func TestMetrics_UpdateGauge(t *testing.T) {
	type fields struct {
		ID    string
		MType string
		Delta *int64
		Value *float64
	}
	type args struct {
		value *float64
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantNewValue *float64
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric := &Metrics{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
			}
			gotNewValue, err := metric.UpdateGauge(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateGauge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNewValue, tt.wantNewValue) {
				t.Errorf("UpdateGauge() gotNewValue = %v, want %v", gotNewValue, tt.wantNewValue)
			}
		})
	}
}
