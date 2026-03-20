package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MapMinder/mapminder_backend/feature/user/domain"
	"github.com/MapMinder/mapminder_backend/feature/user/handler/mock"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	apperror "github.com/MapMinder/mapminder_backend/shared/appError"
	"github.com/MapMinder/mapminder_backend/shared/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestUserHandler_GetUser(t *testing.T) {
	logger.InitForTest()

	testUserId := "test-user-id"
	testUsername := "testuser"
	testEmail := "test@example.com"

	testUser := domain.User{
		UserId:   testUserId,
		Username: testUsername,
		Email:    testEmail,
	}

	tests := []struct {
		name           string
		mockSetup      func(*mock.MockUserUsecase)
		expectedStatus int
	}{
		{
			name: "success",
			mockSetup: func(muu *mock.MockUserUsecase) {
				muu.EXPECT().GetUser(gomock.Any()).Return(testUser, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "failure when usecase returns internal error",
			mockSetup: func(muu *mock.MockUserUsecase) {
				muu.EXPECT().GetUser(gomock.Any()).Return(domain.User{}, apperror.Internal())
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mock.NewMockUserUsecase(ctrl)
			tt.mockSetup(mockUsecase)

			router := gin.New()
			router.Use(middleware.ErrorHandler())

			api := router.Group("/user")
			handler := NewUserHandler(mockUsecase)
			api.Use(func(c *gin.Context) {
				ctx := context.WithValue(c.Request.Context(), middleware.UserIDKey, testUserId)
				c.Request = c.Request.WithContext(ctx)
				c.Next()
			})

			req := httptest.NewRequest(
				http.MethodGet,
				"/user/me",
				nil,
			)

			handler.RegisterRoutes(api)
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