// Package monitor contains bulk of working code for metrics agent service.
package monitor

import (
	"context"
	"fmt"
	"time"

	errorHandling "github.com/evgenytr/metrics.git/internal/errors"
)

// Task struct describes queued task.
type Task struct {
	id int
}

// Queue struct describes task queue.
type Queue struct {
	ch chan *Task
}

// Worker struct describes worker instance.
type Worker struct {
	queue *Queue
	id    int64
}

type workerFunc func() error

// NewQueue creates and returns pointer to new Queue of designated size
func NewQueue(bufSize int64) *Queue {
	if bufSize > 0 {
		return &Queue{
			ch: make(chan *Task, bufSize),
		}
	}

	return &Queue{
		ch: make(chan *Task, 1),
	}
}

func (q *Queue) push(t *Task) {
	q.ch <- t
}

func (q *Queue) pop() *Task {
	return <-q.ch
}

func (q *Queue) close() {
	close(q.ch)
}

// ScheduleTasks fills the Queue with tasks.
func (q *Queue) ScheduleTasks(ctx context.Context, interval time.Duration) {
	taskID := 0
	defer q.close()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("queue stopped scheduling new tasks")
			return
		default:
			time.Sleep(interval)
			q.push(&Task{id: taskID})
			taskID++
		}

	}
}

// NewWorker creates and returns pointer to new Worker instance
func NewWorker(id int64, queue *Queue) *Worker {
	fmt.Println("worker created", id)
	return &Worker{
		id:    id,
		queue: queue,
	}
}

// Loop method sets function to run in worker
func (w *Worker) Loop(ctx context.Context, cancelCtx context.CancelCauseFunc, fn workerFunc) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("worker %w cancelled", w.id)
			return
		default:
		}
		_ = w.queue.pop()

		err := fn()

		if err != nil {
			for _, retryInterval := range errorHandling.RepeatedAttemptsIntervals {
				fmt.Printf("worker %d retrying in %s %s\n", w.id, retryInterval, err)
				time.Sleep(retryInterval)

				err = fn()

				if err == nil {
					break
				}
			}
		}

		if err != nil {
			fmt.Println("worker %w err %w", w.id, err)
			err = fmt.Errorf("worker %d failed: %w", w.id, err)
			cancelCtx(err)

			return
		}

	}
}
