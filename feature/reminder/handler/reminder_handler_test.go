package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MapMinder/mapminder_backend/feature/reminder/domain"
	"github.com/MapMinder/mapminder_backend/feature/reminder/dto"
	"github.com/MapMinder/mapminder_backend/feature/reminder/handler/mock"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	"github.com/MapMinder/mapminder_backend/shared/middleware"
	"github.com/MapMinder/mapminder_backend/shared/validator"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestReminderHandler_CreateReminder(t *testing.T) {
	logger.InitForTest()
	validator.Init()

	testUserId := "test-user-id"
	testReminderId := "test-reminder-uuid"
	testTitle := "test title"
	testDescription := "test description"
	testLatitude := 20.00
	testLongitude := 20.00

	testReminder := dto.Reminder{
		Title:       testTitle,
		Description: testDescription,
		Latitude:    testLatitude,
		Longitude:   testLongitude,
	}

	testReminderRes := domain.Reminder{
		ReminderId:  testReminderId,
		UserId:      testUserId,
		Title:       testTitle,
		Description: testDescription,
		Latitude:    testLatitude,
		Longitude:   testLongitude,
		Radius:      domain.DefaultRadius,
		Status:      string(domain.CreatedStatus),
	}
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(*mock.MockReminderUsecase)
		expectedStatus int
	}{
		{
			name:        "success",
			requestBody: `{"title":"test title","description":"test description","latitude":20.00,"longitude":20.00}`,
			mockSetup: func(mru *mock.MockReminderUsecase) {
				mru.EXPECT().CreateReminder(gomock.Any(), testReminder).Return(testReminderRes, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "failure when latitude is invalid",
			requestBody:    `{"title":"test title","description":"test description","latitude":-200.00,"longitude":20.00}`,
			mockSetup:      func(mru *mock.MockReminderUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "failure when longitude is invalid",
			requestBody:    `{"title":"test title","description":"test description","latitude":20.00,"longitude":-200.00}`,
			mockSetup:      func(mru *mock.MockReminderUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "failure when title is not provided",
			requestBody: `{"description":"test description","latitude":20.00,"longitude":-200.00}`,
			mockSetup: func(mru *mock.MockReminderUsecase) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "failure when description is not provided",
			requestBody:    `{"title":"test title","latitude":20.00,"longitude":-200.00}`,
			mockSetup:      func(mru *mock.MockReminderUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "failure when latitude is not provided",
			requestBody:    `{"title":"test title","description":"test description","longitude":-200.00}`,
			mockSetup:      func(mru *mock.MockReminderUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "failure when longitude is not provided",
			requestBody:    `{"title":"test title","description":"test description","latitude":20.00}`,
			mockSetup:      func(mru *mock.MockReminderUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "failure when falied to bind request",
			requestBody:    ``,
			mockSetup:      func(mru *mock.MockReminderUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mock.NewMockReminderUsecase(ctrl)
			tt.mockSetup(mockUsecase)

			router := gin.New()
			router.Use(middleware.ErrorHandler())

			api := router.Group("/reminder")
			handler := NewReminderHandler(mockUsecase)
			api.Use(func(c *gin.Context) {
				ctx := context.WithValue(c.Request.Context(), middleware.UserIDKey, testUserId)
				c.Request = c.Request.WithContext(ctx)
				c.Next()
			})

			req := httptest.NewRequest(
				http.MethodPost,
				"/reminder/",
				bytes.NewBufferString(tt.requestBody),
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

func TestReminderHandler_GetReminder(t *testing.T) {
	logger.InitForTest()
	validator.Init()

	testUserId := "test-user-id"
	testReminderId := "550e8400-e29b-41d4-a716-446655440000"
	testTitle := "test title"
	testDescription := "test description"
	testLatitude := 20.00
	testLongitude := 20.00

	testReminderRes := domain.Reminder{
		ReminderId:  testReminderId,
		UserId:      testUserId,
		Title:       testTitle,
		Description: testDescription,
		Latitude:    testLatitude,
		Longitude:   testLongitude,
		Radius:      domain.DefaultRadius,
		Status:      string(domain.CreatedStatus),
	}

	tests := []struct {
		name           string
		reminderId     string
		mockSetup      func(*mock.MockReminderUsecase)
		expectedStatus int
	}{
		{
			name:       "success",
			reminderId: testReminderId,
			mockSetup: func(mru *mock.MockReminderUsecase) {
				mru.EXPECT().GetReminder(gomock.Any(), testReminderId).Return(testReminderRes, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:       "failure when reminder id is invalid",
			reminderId: "invalid-uuid",
			mockSetup: func(mru *mock.MockReminderUsecase) {
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mock.NewMockReminderUsecase(ctrl)
			tt.mockSetup(mockUsecase)

			router := gin.New()
			router.Use(middleware.ErrorHandler())

			api := router.Group("/reminder")
			handler := NewReminderHandler(mockUsecase)
			api.Use(func(c *gin.Context) {
				ctx := context.WithValue(c.Request.Context(), middleware.UserIDKey, testUserId)
				c.Request = c.Request.WithContext(ctx)
				c.Next()
			})

			req := httptest.NewRequest(
				http.MethodGet,
				"/reminder/"+tt.reminderId,
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
