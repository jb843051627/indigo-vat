package clock

import (
	"context"
	"sync"
	"time"
)

type Fixed struct {
	mu    sync.RWMutex
	value time.Time
}

func NewFixed(value time.Time) *Fixed { return &Fixed{value: value} }
func (f *Fixed) Now() time.Time       { f.mu.RLock(); defer f.mu.RUnlock(); return f.value }
func (f *Fixed) Advance(delay time.Duration) {
	f.mu.Lock()
	f.value = f.value.Add(delay)
	f.mu.Unlock()
}
func (f *Fixed) Sleep(ctx context.Context, delay time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		f.Advance(delay)
		return nil
	}
}
