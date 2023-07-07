package monitor

import (
	"context"
	"fmt"
	"time"
)

type Task struct {
	id int
}

type Queue struct {
	ch chan *Task
}

func NewQueue(bufSize *int64) *Queue {
	if bufSize != nil && *bufSize > 0 {
		return &Queue{
			ch: make(chan *Task, *bufSize),
		}
	}

	return &Queue{
		ch: make(chan *Task, 1),
	}
}

func (q *Queue) Push(t *Task) {
	q.ch <- t
}

func (q *Queue) Pop() *Task {
	return <-q.ch
}

func (q *Queue) ScheduleTasks(interval *time.Duration) {
	taskID := 0
	for {
		time.Sleep(*interval)
		q.Push(&Task{id: taskID})
		taskID++
	}
}

type Worker struct {
	id    int64
	queue *Queue
}

func NewWorker(id int64, queue *Queue) *Worker {
	return &Worker{
		id:    id,
		queue: queue,
	}
}

type workerFunc func() error

func (w *Worker) Loop(ctx context.Context, fn workerFunc) {
	_, cancelCtx := context.WithCancelCause(ctx)
	for {
		_ = w.queue.Pop()

		err := fn()

		if err != nil {
			fmt.Println("worker %w err %w", w.id, err)
			cancelCtx(err)
			return
		}

	}
}
