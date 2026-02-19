package handler

import (
	"net/http"

	"github.com/MapMinder/mapminder_backend/feature/auth/dto"
	"github.com/MapMinder/mapminder_backend/feature/auth/usecase"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	"github.com/MapMinder/mapminder_backend/internal/status"
	apperror "github.com/MapMinder/mapminder_backend/shared/appError"
	"github.com/MapMinder/mapminder_backend/shared/validator"
	"github.com/gin-gonic/gin"
)

// GoogleSignInHandler
type GoogleSignInHandler struct {
	GoogleSignInUsecase usecase.GoogleSignInUsecase
}

// NewGoogleSignInHandler
func NewGoogleSignInhandler(googleSignInUsecase usecase.GoogleSignInUsecase) *GoogleSignInHandler {
	return &GoogleSignInHandler{
		GoogleSignInUsecase: googleSignInUsecase,
	}
}

// RegisterRoutes add google endpoints
func (h *GoogleSignInHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/google", h.GoogleSignIn)
}

// GoogleSignIn login user using google sign in
func (h *GoogleSignInHandler) GoogleSignIn(r *gin.Context) {
	logger.Info("GoogleSignIn")
	ctx := r.Request.Context()

	var req dto.GoogleSignInToken
	if err := r.ShouldBind(&req); err != nil {
		logger.Errorw("Failed To Bind Request", err)
		r.Error(apperror.BadRequest())
		return
	}

	if err := validator.ValidateStruct(req); err != nil {
		logger.Errorw("Failed To Validate Request", err)
		r.Error(apperror.BadRequest())
		return
	}

	jwtToken, err := h.GoogleSignInUsecase.GoogleSignIn(ctx, req.IdToken)
	if err != nil {
		logger.Error(err)
		r.Error(err)
		return
	}

	res := dto.GoogleSignInRes{
		Status: status.Success,
		Result: jwtToken,
	}

	r.JSON(http.StatusOK, res)
}
