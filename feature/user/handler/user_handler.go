package handler

import (
	"net/http"

	"github.com/MapMinder/mapminder_backend/feature/user/dto"
	"github.com/MapMinder/mapminder_backend/feature/user/mapper"
	"github.com/MapMinder/mapminder_backend/feature/user/usecase"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	"github.com/MapMinder/mapminder_backend/internal/status"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		UserUsecase: userUsecase,
	}
}

func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/me", h.GetUser)
}

func (h *UserHandler) GetUser(r *gin.Context) {
	logger.Infof("user handler: GetUser")
	ctx := r.Request.Context()

	user, err := h.UserUsecase.GetUser(ctx)
	if err != nil {
		r.Error(err)
		return
	}

	res := dto.UserRes{
		Status: status.Success,
		Result: mapper.MapUser(user),
	}

	r.JSON(http.StatusOK, res)
}
