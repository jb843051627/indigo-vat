package worker

import (
	"context"
	"github.com/jb843051627/indigo-vat/internal/model"
	"github.com/jb843051627/indigo-vat/internal/service"
	"time"
)

type Maturity struct {
	service  *service.Service
	interval time.Duration
}

func NewMaturity(s *service.Service, interval time.Duration) *Maturity {
	if interval < time.Millisecond {
		interval = time.Second
	}
	return &Maturity{service: s, interval: interval}
}
func (m *Maturity) Run(ctx context.Context) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.scan(ctx)
		}
	}
}
func (m *Maturity) scan(ctx context.Context) {
	cycles, err := m.service.ListCycles(ctx, model.CycleFermenting)
	if err != nil {
		return
	}
	for _, cycle := range cycles {
		if cycle.StartedAt.Add(time.Hour).Before(time.Now()) {
			_, _ = m.service.AdvanceCycle(ctx, cycle.ID, model.CycleMatured, cycle.Revision, "maturity scan")
		}
	}
}
