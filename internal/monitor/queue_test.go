package monitor

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestNewQueue(t *testing.T) {
	type args struct {
		bufSize int64
	}
	tests := []struct {
		name string
		args args
		want *Queue
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueue(tt.args.bufSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWorker(t *testing.T) {
	type args struct {
		id    int64
		queue *Queue
	}
	tests := []struct {
		name string
		args args
		want *Worker
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWorker(tt.args.id, tt.args.queue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWorker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueue_ScheduleTasks(t *testing.T) {
	type fields struct {
		ch chan *Task
	}
	type args struct {
		interval time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queue{
				ch: tt.fields.ch,
			}
			q.ScheduleTasks(tt.args.interval)
		})
	}
}

func TestQueue_close(t *testing.T) {
	type fields struct {
		ch chan *Task
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queue{
				ch: tt.fields.ch,
			}
			q.close()
		})
	}
}

func TestQueue_pop(t *testing.T) {
	type fields struct {
		ch chan *Task
	}
	tests := []struct {
		name   string
		fields fields
		want   *Task
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queue{
				ch: tt.fields.ch,
			}
			if got := q.pop(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueue_push(t *testing.T) {
	type fields struct {
		ch chan *Task
	}
	type args struct {
		t *Task
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queue{
				ch: tt.fields.ch,
			}
			q.push(tt.args.t)
		})
	}
}

func TestWorker_Loop(t *testing.T) {
	type fields struct {
		id    int64
		queue *Queue
	}
	type args struct {
		ctx       context.Context
		cancelCtx context.CancelCauseFunc
		fn        workerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Worker{
				id:    tt.fields.id,
				queue: tt.fields.queue,
			}
			w.Loop(tt.args.ctx, tt.args.cancelCtx, tt.args.fn)
		})
	}
}