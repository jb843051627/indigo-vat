package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/indigo-vat/internal/model"
	"github.com/jb843051627/indigo-vat/internal/store"
	"github.com/jb843051627/indigo-vat/internal/validation"
)

func (s *Service) CreateVat(ctx context.Context, in model.VatInput) (model.Vat, error) {
	if err := validation.ValidateVat(in); err != nil {
		return model.Vat{}, err
	}
	now := s.clock.Now()
	v := model.Vat{ID: validation.NewID("vat"), Name: validation.CleanText(in.Name), Site: validation.CleanText(in.Site), State: model.VatActive, Capacity: in.Capacity, TargetPH: in.TargetPH, TargetTemp: in.TargetTemp, Timezone: validation.FirstNonEmpty(in.Timezone, "UTC"), CreatedAt: now, UpdatedAt: now}
	if err := s.db.PutVat(ctx, v); err != nil {
		return model.Vat{}, err
	}
	if err := s.db.AppendAudit(ctx, audit("vat", v.ID, "created", v.Name, now)); err != nil {
		return model.Vat{}, err
	}
	return v, nil
}
func (s *Service) GetVat(ctx context.Context, id string) (model.Vat, error) {
	if err := ensureID(id); err != nil {
		return model.Vat{}, err
	}
	v, err := s.db.GetVat(ctx, id)
	if err != nil {
		return model.Vat{}, wrapNotFound("vat "+id, err)
	}
	return v, nil
}
func (s *Service) ListVats(ctx context.Context) ([]model.Vat, error) { return s.db.ListVats(ctx) }
func (s *Service) changeVat(ctx context.Context, id, state string) (model.Vat, error) {
	v, err := s.GetVat(ctx, id)
	if err != nil {
		return model.Vat{}, err
	}
	if v.State == model.VatRetired {
		return model.Vat{}, fmt.Errorf("%w: retired", ErrConflict)
	}
	when := s.nowText()
	if err := s.db.SetVatState(context.Background(), id, state, when); err != nil {
		return model.Vat{}, wrapNotFound("vat "+id, err)
	}
	v.State = state
	return v, s.db.AppendAudit(ctx, audit("vat", id, "state_changed", state, s.clock.Now()))
}
func (s *Service) PauseVat(ctx context.Context, id string) (model.Vat, error) {
	return s.changeVat(ctx, id, model.VatPaused)
}
func (s *Service) ResumeVat(ctx context.Context, id string) (model.Vat, error) {
	return s.changeVat(ctx, id, model.VatActive)
}
func (s *Service) RetireVat(ctx context.Context, id string) (model.Vat, error) {
	return s.changeVat(ctx, id, model.VatRetired)
}

var _ *store.DB
