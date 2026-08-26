package clock

import (
	"context"
	"time"
)

type Clock interface {
	Now() time.Time
	Sleep(context.Context, time.Duration) error
}
type Real struct{}

func (Real) Now() time.Time { return time.Now().UTC() }
func (Real) Sleep(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
