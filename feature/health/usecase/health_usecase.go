package usecase

import (
	"context"

	"github.com/MapMinder/mapminder_backend/feature/health/repository"
)

type HealthUsecase interface {
	CheckHealth(ctx context.Context) (*HealthStatus, error)
}

type HealthStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type healthUsecase struct {
	healthRepo repository.HealthRepository
}

func NewHealthUsecase(healthRepo repository.HealthRepository) HealthUsecase {
	return &healthUsecase{
		healthRepo: healthRepo,
	}
}
func (u *healthUsecase) CheckHealth(ctx context.Context) (*HealthStatus, error) {

	return &HealthStatus{
		Status: "ok",
	}, nil
}
