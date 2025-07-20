package handler

import (
	"net/http"

	"github.com/MapMinder/mapminder_backend/feature/health/usecase"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	healthUsecase usecase.HealthUsecase
}

func NewHealthHandler(healthUsecase usecase.HealthUsecase) *HealthHandler {
	return &HealthHandler{
		healthUsecase: healthUsecase,
	}
}

func (h *HealthHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", h.CheckHealth)
}

func (h *HealthHandler) CheckHealth(c *gin.Context) {
	status, err := h.healthUsecase.CheckHealth(c)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, status)
}
