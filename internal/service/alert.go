package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jb843051627/indigo-vat/internal/model"
	"github.com/jb843051627/indigo-vat/internal/validation"
)

func (s *Service) RaiseAlert(ctx context.Context, in model.AlertInput) (model.Alert, error) {
	if in.CycleID == "" || in.Code == "" {
		return model.Alert{}, fmt.Errorf("%w: alert", validation.ErrInvalidInput)
	}
	a := model.Alert{ID: validation.NewID("alert"), CycleID: in.CycleID, Level: in.Level, Code: in.Code, Message: in.Message, State: model.AlertOpen, CreatedAt: s.clock.Now()}
	return a, s.db.PutAlert(ctx, a)
}
func (s *Service) ListAlerts(ctx context.Context, cycleID string) ([]model.Alert, error) {
	return s.db.ListAlerts(ctx, cycleID)
}
func (s *Service) AcknowledgeAlert(ctx context.Context, id string) (model.Alert, error) {
	a, err := s.db.ListAlerts(ctx, "")
	if err != nil {
		return model.Alert{}, err
	}
	for _, item := range a {
		if item.ID == id {
			item.State = model.AlertAcknowledged
			item.AcknowledgedAt = s.clock.Now()
			if err := s.db.AcknowledgeAlert(ctx, id, storeTime(item.AcknowledgedAt)); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return model.Alert{}, wrapNotFound("alert "+id, err)
				}
				return model.Alert{}, err
			}
			return item, nil
		}
	}
	return model.Alert{}, ErrNotFound
}
