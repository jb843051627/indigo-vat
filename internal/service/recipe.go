package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/indigo-vat/internal/model"
	"github.com/jb843051627/indigo-vat/internal/validation"
	"time"
)

func (s *Service) CreateRecipe(ctx context.Context, in model.RecipeInput) (model.Recipe, error) {
	if err := validation.ValidateRecipe(in); err != nil {
		return model.Recipe{}, err
	}
	now := s.clock.Now()
	r := model.Recipe{ID: validation.NewID("recipe"), Name: validation.CleanText(in.Name), TargetHue: in.TargetHue, HueTolerance: in.HueTolerance, MinMinutes: in.MinMinutes, MaxMinutes: in.MaxMinutes, State: model.RecipeDraft, CreatedAt: now, UpdatedAt: now}
	if err := s.db.PutRecipe(ctx, r); err != nil {
		return model.Recipe{}, err
	}
	return r, s.db.AppendAudit(ctx, audit("recipe", r.ID, "created", r.Name, now))
}
func (s *Service) GetRecipe(ctx context.Context, id string) (model.Recipe, error) {
	r, err := s.db.GetRecipe(ctx, id)
	if err != nil {
		return model.Recipe{}, wrapNotFound("recipe "+id, err)
	}
	return r, nil
}
func (s *Service) AddStage(ctx context.Context, id string, in model.StageInput) (model.RecipeStage, error) {
	if err := validation.ValidateStage(in); err != nil {
		return model.RecipeStage{}, err
	}
	r, err := s.GetRecipe(ctx, id)
	if err != nil {
		return model.RecipeStage{}, err
	}
	if r.State != model.RecipeDraft {
		return model.RecipeStage{}, ErrConflict
	}
	stage := model.RecipeStage{ID: validation.NewID("stage"), RecipeID: id, Position: len(r.Stages) + 1, Name: validation.CleanText(in.Name), Minutes: in.Minutes, PHMin: in.PHMin, PHMax: in.PHMax, TempMin: in.TempMin, TempMax: in.TempMax}
	r.Stages = append(r.Stages, stage)
	if err := s.db.PutRecipe(ctx, r); err != nil {
		return model.RecipeStage{}, err
	}
	return stage, nil
}
func (s *Service) PublishRecipe(ctx context.Context, id string) (model.Recipe, error) {
	r, err := s.GetRecipe(ctx, id)
	if err != nil {
		return model.Recipe{}, err
	}
	if len(r.Stages) == 0 {
		return model.Recipe{}, fmt.Errorf("%w: stages", ErrNotReady)
	}
	now := s.clock.Now()
	if err := s.db.SetRecipeState(ctx, id, model.RecipeReady, now.Format(time.RFC3339Nano)); err != nil {
		return model.Recipe{}, err
	}
	r.State = model.RecipeReady
	return r, nil
}
func (s *Service) ListRecipes(ctx context.Context) ([]model.Recipe, error) { return nil, nil }
