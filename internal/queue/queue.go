package queue

import (
	"context"
	"errors"
	"sync"
)

var ErrClosed = errors.New("queue is closed")

type Job func(context.Context) error
type Queue struct {
	jobs   chan Job
	done   chan struct{}
	mu     sync.RWMutex
	closed bool
	once   sync.Once
	wg     sync.WaitGroup
}

func New(size int) *Queue {
	if size < 1 {
		size = 1
	}
	return &Queue{jobs: make(chan Job, size), done: make(chan struct{})}
}
func (q *Queue) Submit(ctx context.Context, job Job) error {
	q.mu.RLock()
	closed := q.closed
	q.mu.RUnlock()
	if closed {
		return ErrClosed
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-q.done:
		return ErrClosed
	case q.jobs <- job:
		return nil
	}
}
func (q *Queue) Start(ctx context.Context, workers int, handle func(context.Context, Job) error) {
	if workers < 1 {
		workers = 1
	}
	q.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go q.worker(ctx, handle)
	}
}
func (q *Queue) worker(ctx context.Context, handle func(context.Context, Job) error) {
	defer q.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case <-q.done:
			return
		case job, ok := <-q.jobs:
			if !ok {
				return
			}
			_ = handle(ctx, job)
		}
	}
}
func (q *Queue) Close() {
	q.once.Do(func() {
		q.mu.Lock()
		q.closed = true
		q.mu.Unlock()
		close(q.done)
		q.wg.Wait()
	})
}
func (q *Queue) Done() <-chan struct{} { return q.done }
