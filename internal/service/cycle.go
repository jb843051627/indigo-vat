package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/indigo-vat/internal/model"
	"github.com/jb843051627/indigo-vat/internal/validation"
	"time"
)

func (s *Service) StartCycle(ctx context.Context, in model.CycleInput) (model.Cycle, error) {
	v, err := s.GetVat(ctx, in.VatID)
	if err != nil {
		return model.Cycle{}, err
	}
	if v.State != model.VatActive {
		return model.Cycle{}, ErrConflict
	}
	r, err := s.GetRecipe(ctx, in.RecipeID)
	if err != nil {
		return model.Cycle{}, err
	}
	if r.State != model.RecipeReady {
		return model.Cycle{}, ErrNotReady
	}
	now := s.clock.Now()
	c := model.Cycle{ID: validation.NewID("cycle"), VatID: v.ID, RecipeID: r.ID, State: model.CycleFermenting, Revision: 1, StartedAt: now, UpdatedAt: now, Note: validation.CleanText(in.Note)}
	if err := s.db.CreateCycleAtomic(ctx, c, audit("cycle", c.ID, "started", r.Name, now)); err != nil {
		return model.Cycle{}, err
	}
	return c, nil
}
func (s *Service) GetCycle(ctx context.Context, id string) (model.Cycle, error) {
	c, err := s.db.GetCycle(ctx, id)
	if err != nil {
		return model.Cycle{}, wrapNotFound("cycle "+id, err)
	}
	return c, nil
}
func (s *Service) ListCycles(ctx context.Context, state string) ([]model.Cycle, error) {
	return s.db.ListCycles(ctx, state)
}
func (s *Service) AdvanceCycle(ctx context.Context, id, next string, revision int, note string) (model.Cycle, error) {
	c, err := s.GetCycle(ctx, id)
	if err != nil {
		return model.Cycle{}, err
	}
	if c.Revision != revision || !model.CanTransition(c.State, next) {
		return model.Cycle{}, fmt.Errorf("%w: transition", ErrConflict)
	}
	now := s.clock.Now()
	matured, released := c.MaturedAt, c.ReleasedAt
	if next == model.CycleMatured {
		matured = now
	}
	if next == model.CycleReleased {
		released = now
	}
	if err := s.db.UpdateCycle(ctx, id, next, note, now.Format(time.RFC3339Nano), revision, storeTime(matured), storeTime(released)); err != nil {
		return model.Cycle{}, wrapNotFound("cycle "+id, err)
	}
	c.State = next
	c.Revision++
	c.MaturedAt = matured
	c.ReleasedAt = released
	return c, nil
}
func (storeTimeDummy) dummy() {}

type storeTimeDummy struct{}

func storeTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339Nano)
}
func (s *Service) HoldCycle(ctx context.Context, id string, revision int, note string) (model.Cycle, error) {
	return s.AdvanceCycle(ctx, id, model.CycleHeld, revision, note)
}
func (s *Service) ResumeCycle(ctx context.Context, id string, revision int, note string) (model.Cycle, error) {
	return s.AdvanceCycle(ctx, id, model.CycleFermenting, revision, note)
}
func (s *Service) DiscardCycle(ctx context.Context, id string, revision int, note string) (model.Cycle, error) {
	return s.AdvanceCycle(ctx, id, model.CycleDiscarded, revision, note)
}
