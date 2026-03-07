package usecase

import (
	"context"

	"github.com/MapMinder/mapminder_backend/feature/health/repository"
	"github.com/MapMinder/mapminder_backend/shared/response"
)

type HealthUsecase interface {
	CheckHealth(ctx context.Context) (*response.APIResponse, error)
}

type healthUsecase struct {
	healthRepo repository.HealthRepository
}

func NewHealthUsecase(healthRepo repository.HealthRepository) HealthUsecase {
	return &healthUsecase{
		healthRepo: healthRepo,
	}
}

func (u *healthUsecase) CheckHealth(ctx context.Context) (*response.APIResponse, error) {
	return &response.APIResponse{
		Status: "ok",
	}, nil
}
