package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jb843051627/indigo-vat/internal/clock"
	"github.com/jb843051627/indigo-vat/internal/metrics"
	"github.com/jb843051627/indigo-vat/internal/model"
	"github.com/jb843051627/indigo-vat/internal/queue"
	"github.com/jb843051627/indigo-vat/internal/store"
	"github.com/jb843051627/indigo-vat/internal/validation"
	"sync"
	"time"
)

var ErrNotFound = errors.New("entity not found")
var ErrConflict = errors.New("state conflict")
var ErrNotReady = errors.New("entity not ready")

type Service struct {
	db      *store.DB
	clock   clock.Clock
	queue   *queue.Queue
	metrics *metrics.Registry
	startMu sync.Mutex
	started bool
}

func New(db *store.DB) *Service { return NewWithClock(db, clock.Real{}) }
func NewWithClock(db *store.DB, c clock.Clock) *Service {
	return &Service{db: db, clock: c, queue: queue.New(32), metrics: metrics.New()}
}
func (s *Service) Start(ctx context.Context) {
	s.startMu.Lock()
	defer s.startMu.Unlock()
	if s.started {
		return
	}
	s.started = true
	s.queue.Start(ctx, 2, func(jobCtx context.Context, job queue.Job) error {
		s.metrics.Inc("jobs.started")
		if err := job(jobCtx); err != nil {
			s.metrics.Inc("jobs.failed")
			return err
		}
		s.metrics.Inc("jobs.completed")
		return nil
	})
}
func (s *Service) Close() {
	s.startMu.Lock()
	started := s.started
	s.startMu.Unlock()
	if !started {
		return
	}
	s.queue.Close()
}
func (s *Service) DB() *store.DB { return s.db }
func (s *Service) Metrics() map[string]int64 {
	return s.metrics.Snapshot()
}
func (s *Service) nowText() string { return s.clock.Now().UTC().Format(time.RFC3339Nano) }
func audit(typ, id, action, detail string, now time.Time) model.AuditEvent {
	return model.AuditEvent{ID: validation.NewID("audit"), EntityType: typ, EntityID: id, Action: action, Detail: detail, CreatedAt: now}
}
func wrapNotFound(name string, err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	return err
}
func ensureID(id string) error {
	if !validation.IsID(id) {
		return fmt.Errorf("%w: bad id", validation.ErrInvalidInput)
	}
	return nil
}

func checkContext(ctx context.Context) error { return ctx.Err() }

func conflictError(detail string) error { return fmt.Errorf("%w: %s", ErrConflict, detail) }

func wrapSampleValidation(err error) error { return fmt.Errorf("%w: sample validation", err) }
