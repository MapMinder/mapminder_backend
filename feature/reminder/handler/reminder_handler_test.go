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
	apperror "github.com/MapMinder/mapminder_backend/shared/appError"
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
		Status:      string(domain.ActiveStatus),
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
				"/reminder",
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
		Status:      string(domain.ActiveStatus),
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

func TestReminderHandler_GetReminders(t *testing.T) {
	logger.InitForTest()
	validator.Init()

	testUserId := "test-user-id"
	testReminderId := "550e8400-e29b-41d4-a716-446655440000"
	testTitle := "test title"
	testDescription := "test description"
	testLatitude := 20.00
	testLongitude := 20.00

	testReminderRes := []domain.Reminder{
		{
			ReminderId:  testReminderId,
			UserId:      testUserId,
			Title:       testTitle,
			Description: testDescription,
			Latitude:    testLatitude,
			Longitude:   testLongitude,
			Radius:      domain.DefaultRadius,
			Status:      string(domain.ActiveStatus),
		},
	}

	tests := []struct {
		name           string
		status         string
		mockSetup      func(*mock.MockReminderUsecase)
		expectedStatus int
	}{
		{
			name:   "success",
			status: "status=active",
			mockSetup: func(mru *mock.MockReminderUsecase) {
				mru.EXPECT().GetReminders(gomock.Any(), "active").Return(testReminderRes, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "fail when invalid status is provided",
			status: "status=invalid-status",
			mockSetup: func(mru *mock.MockReminderUsecase) {
				mru.EXPECT().GetReminders(gomock.Any(), "invalid-status").Return([]domain.Reminder{}, apperror.BadRequest())
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "fails with internal error",
			status: "status=active",
			mockSetup: func(mru *mock.MockReminderUsecase) {
				mru.EXPECT().GetReminders(gomock.Any(), "active").Return([]domain.Reminder{}, apperror.Internal())
			},
			expectedStatus: http.StatusInternalServerError,
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
				"/reminder?"+tt.status,
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

func TestReminderHandler_DeleteReminder(t *testing.T) {
	logger.InitForTest()
	validator.Init()

	testUserId := "test-user-id"
	testReminderId := "550e8400-e29b-41d4-a716-446655440000"

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
				mru.EXPECT().DeleteReminder(gomock.Any(), testReminderId).Return(nil)
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
				http.MethodDelete,
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

func TestReminderHandler_UpdateReminder(t *testing.T) {
	logger.InitForTest()
	validator.Init()

	testUserId := "test-user-id"
	testReminderId := "550e8400-e29b-41d4-a716-446655440000"
	testTitle := "test title"
	testDescription := "test description"
	testLatitude := 20.00
	testLongitude := 20.00

	testReminder := dto.UpdateReminder{
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
		Status:      string(domain.ActiveStatus),
	}
	tests := []struct {
		name           string
		reminderId     string
		requestBody    string
		mockSetup      func(*mock.MockReminderUsecase)
		expectedStatus int
	}{
		{
			name:        "success",
			reminderId:  testReminderId,
			requestBody: `{"title":"test title","description":"test description","latitude":20.00,"longitude":20.00}`,
			mockSetup: func(mru *mock.MockReminderUsecase) {
				mru.EXPECT().UpdateReminder(gomock.Any(), testReminderId, testReminder).Return(testReminderRes, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "success when title is not given",
			reminderId:  testReminderId,
			requestBody: `{"description":"test description","latitude":20.00,"longitude":20.00}`,
			mockSetup: func(mru *mock.MockReminderUsecase) {
				mru.EXPECT().UpdateReminder(gomock.Any(), testReminderId, dto.UpdateReminder{
					Description: testDescription,
					Latitude:    testLatitude,
					Longitude:   testLongitude,
				}).Return(testReminderRes, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "success when description is not given",
			reminderId:  testReminderId,
			requestBody: `{"title":"test title","latitude":20.00,"longitude":20.00}`,
			mockSetup: func(mru *mock.MockReminderUsecase) {
				mru.EXPECT().UpdateReminder(gomock.Any(), testReminderId, dto.UpdateReminder{
					Title:     testTitle,
					Latitude:  testLatitude,
					Longitude: testLongitude,
				}).Return(testReminderRes, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "success when latitude is not given",
			reminderId:  testReminderId,
			requestBody: `{"title":"test title","description":"test description","longitude":20.00}`,
			mockSetup: func(mru *mock.MockReminderUsecase) {
				mru.EXPECT().UpdateReminder(gomock.Any(), testReminderId, dto.UpdateReminder{
					Title:       testTitle,
					Description: testDescription,
					Longitude:   testLongitude,
				}).Return(testReminderRes, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "success when longitude is not given",
			reminderId:  testReminderId,
			requestBody: `{"title":"test title","description":"test description","latitude":20.00}`,
			mockSetup: func(mru *mock.MockReminderUsecase) {
				mru.EXPECT().UpdateReminder(gomock.Any(), testReminderId, dto.UpdateReminder{
					Title:       testTitle,
					Description: testDescription,
					Latitude:    testLatitude,
				}).Return(testReminderRes, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "failure when reminder id is invalid",
			reminderId:     "invalid-uuid",
			requestBody:    `{"title":"test title","description":"test description","latitude":20.00}`,
			mockSetup:      func(mru *mock.MockReminderUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "failure when invalid latitude is given",
			reminderId:     testReminderId,
			requestBody:    `{"title":"test title","description":"test description","latitude":-200.00,"longitude":20.00}`,
			mockSetup:      func(mru *mock.MockReminderUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "failure when invalid longitude is given",
			reminderId:     testReminderId,
			requestBody:    `{"title":"test title","description":"test description","latitude":20.00,"longitude":-200.00}`,
			mockSetup:      func(mru *mock.MockReminderUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "failure when usecase return internal error",
			reminderId:  testReminderId,
			requestBody: `{"title":"test title","description":"test description","latitude":20.00,"longitude":20.00}`,
			mockSetup: func(mru *mock.MockReminderUsecase) {
				mru.EXPECT().UpdateReminder(gomock.Any(), testReminderId, testReminder).Return(domain.Reminder{}, apperror.Internal())
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:        "failure when usecase return unauthorized error",
			reminderId:  testReminderId,
			requestBody: `{"title":"test title","description":"test description","latitude":20.00,"longitude":20.00}`,
			mockSetup: func(mru *mock.MockReminderUsecase) {
				mru.EXPECT().UpdateReminder(gomock.Any(), testReminderId, testReminder).Return(domain.Reminder{}, apperror.Unauthorized())
			},
			expectedStatus: http.StatusUnauthorized,
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
				http.MethodPatch,
				"/reminder/"+tt.reminderId,
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
