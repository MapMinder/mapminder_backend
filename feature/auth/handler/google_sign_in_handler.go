package handler

import (
	"github.com/MapMinder/mapminder_backend/feature/auth/usecase"
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

// RegisterRoutes
func (h *GoogleSignInHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/google", h.GoogleSignIn)
}

func (h *GoogleSignInHandler) GoogleSignIn(r *gin.Context) {
}
