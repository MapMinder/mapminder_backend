package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MapMinder/mapminder_backend/feature/auth/handler/mock"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	"github.com/MapMinder/mapminder_backend/internal/middleware"
	apperror "github.com/MapMinder/mapminder_backend/shared/appError"
	"github.com/MapMinder/mapminder_backend/shared/validator"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestGoogleSignInHandler_GoogleSignIn(t *testing.T) {
	logger.InitForTest()
	validator.Init()
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(*mock.MockGoogleSignInUsecase)
		expectedStatus int
	}{
		{
			name:        "success",
			requestBody: `{"id_token":"valid-token"}`,
			mockSetup: func(m *mock.MockGoogleSignInUsecase) {
				m.EXPECT().
					GoogleSignIn(gomock.Any(), "valid-token").
					Return("jwt-token", nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "failed to validate",
			requestBody:    `{"id_token":""}`,
			mockSetup:      func(m *mock.MockGoogleSignInUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "usecase returns internal error",
			requestBody: `{"id_token":"valid-token"}`,
			mockSetup: func(m *mock.MockGoogleSignInUsecase) {
				m.EXPECT().
					GoogleSignIn(gomock.Any(), "valid-token").
					Return("", apperror.Internal())
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:        "usecase returns internal error",
			requestBody: `{"id_token":"valid-token"}`,
			mockSetup: func(m *mock.MockGoogleSignInUsecase) {
				m.EXPECT().
					GoogleSignIn(gomock.Any(), "valid-token").
					Return("", apperror.Internal())
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:        "google api falis to validate id token",
			requestBody: `{"id_token":"valid-token"}`,
			mockSetup: func(m *mock.MockGoogleSignInUsecase) {
				m.EXPECT().
					GoogleSignIn(gomock.Any(), "valid-token").
					Return("", apperror.GoogleError())
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mock.NewMockGoogleSignInUsecase(ctrl)
			tt.mockSetup(mockUsecase)

			router := gin.New()
			router.Use(middleware.ErrorHandler())

			api := router.Group("/auth")
			handler := NewGoogleSignInhandler(mockUsecase)
			handler.RegisterRoutes(api)

			req := httptest.NewRequest(
				http.MethodPost,
				"/auth/google",
				bytes.NewBufferString(tt.requestBody),
			)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d, body %s",
					tt.expectedStatus,
					w.Code,
					w.Body.String(),
				)
			}
		})
	}
}
