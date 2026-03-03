package usecase

import (
	"context"
	"testing"

	"github.com/MapMinder/mapminder_backend/feature/reminder/domain"
	"github.com/MapMinder/mapminder_backend/feature/reminder/dto"
	"github.com/MapMinder/mapminder_backend/feature/reminder/usecase/mock"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	apperror "github.com/MapMinder/mapminder_backend/shared/appError"
	"github.com/MapMinder/mapminder_backend/shared/middleware"
	mockTime "github.com/MapMinder/mapminder_backend/shared/time/mock"
	txMock "github.com/MapMinder/mapminder_backend/shared/tx/mock"
	mockUUID "github.com/MapMinder/mapminder_backend/shared/uuid_manager/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func mockWithinTransaction(t *testing.T, tx *txMock.MockManager) {
	tx.EXPECT().
		WithinTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		}).
		AnyTimes()
}

func TestCreateReminder(t *testing.T) {
	logger.InitForTest()
	ctx := context.Background()
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
		name        string
		args        dto.Reminder
		prepareFunc func(
			mrr *mock.MockReminderRepository,
			mUUID *mockUUID.MockUUIDManager,
			mTx *txMock.MockManager,
			mTime *mockTime.MockRealTimeProvider,
		)
		wantedError error
		wantedRes   domain.Reminder
	}{
		{
			name: "success",
			args: testReminder,
			prepareFunc: func(
				mrr *mock.MockReminderRepository,
				mUUID *mockUUID.MockUUIDManager,
				mTx *txMock.MockManager,
				mTime *mockTime.MockRealTimeProvider,
			) {
				mUUID.EXPECT().NewV7().Return(testReminderId, nil)
				mrr.EXPECT().Create(gomock.Any(), testReminderRes).Return(nil)
			},
			wantedError: nil,
			wantedRes:   testReminderRes,
		},
		{
			name: "fails when saving reminder",
			args: testReminder,
			prepareFunc: func(
				mrr *mock.MockReminderRepository,
				mUUID *mockUUID.MockUUIDManager,
				mTx *txMock.MockManager,
				mTime *mockTime.MockRealTimeProvider,
			) {
				mUUID.EXPECT().NewV7().Return(testReminderId, nil)
				mrr.EXPECT().Create(gomock.Any(), testReminderRes).Return(apperror.Internal())
			},
			wantedError: apperror.Internal(),
			wantedRes:   domain.Reminder{},
		},
		{
			name: "fails when generating uuid",
			args: testReminder,
			prepareFunc: func(
				mrr *mock.MockReminderRepository,
				mUUID *mockUUID.MockUUIDManager,
				mTx *txMock.MockManager,
				mTime *mockTime.MockRealTimeProvider,
			) {
				mUUID.EXPECT().NewV7().Return("", apperror.Internal())
			},
			wantedError: apperror.Internal(),
			wantedRes:   domain.Reminder{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockReminderRepo := mock.NewMockReminderRepository(ctrl)
			mockTx := txMock.NewMockManager(ctrl)
			mockUUID := mockUUID.NewMockUUIDManager(ctrl)
			mockTime := mockTime.NewMockRealTimeProvider(ctrl)

			mockWithinTransaction(t, mockTx)

			tt.prepareFunc(mockReminderRepo, mockUUID, mockTx, mockTime)
			uc := NewReminderUsecase(mockTx, mockUUID, mockReminderRepo, mockTime)

			ctx = context.WithValue(ctx, middleware.UserIDKey, testUserId)
			actualRes, err := uc.CreateReminder(ctx, tt.args)
			if tt.wantedError != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantedError)
				}
				assert.EqualError(t, err, tt.wantedError.Error())
				return
			}

			assert.Equal(t, tt.wantedRes, actualRes)
		})
	}
}

func TestGetReminder(t *testing.T) {
	logger.InitForTest()
	ctx := context.Background()
	testUserId := "test-user-id"
	testReminderId := "test-reminder-uuid"
	testTitle := "test title"
	testDescription := "test description"
	testLatitude := 20.00
	testLongitude := 20.00

	expectedReminder := domain.Reminder{
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
		name        string
		reminderId  string
		prepareFunc func(
			mrr *mock.MockReminderRepository,
			mUUID *mockUUID.MockUUIDManager,
			mTx *txMock.MockManager,
			mTime *mockTime.MockRealTimeProvider,
		)
		wantedError error
		wantedRes   domain.Reminder
	}{
		{
			name:       "success",
			reminderId: testReminderId,
			prepareFunc: func(
				mrr *mock.MockReminderRepository,
				mUUID *mockUUID.MockUUIDManager,
				mTx *txMock.MockManager,
				mTime *mockTime.MockRealTimeProvider,
			) {
				mrr.EXPECT().GetReminder(gomock.Any(), testReminderId).Return(expectedReminder, nil)
			},
			wantedError: nil,
			wantedRes:   expectedReminder,
		},
		{
			name:       "fails when repository returns error",
			reminderId: testReminderId,
			prepareFunc: func(
				mrr *mock.MockReminderRepository,
				mUUID *mockUUID.MockUUIDManager,
				mTx *txMock.MockManager,
				mTime *mockTime.MockRealTimeProvider,
			) {
				mrr.EXPECT().GetReminder(gomock.Any(), testReminderId).Return(domain.Reminder{}, apperror.Internal())
			},
			wantedError: apperror.Internal(),
			wantedRes:   domain.Reminder{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockReminderRepo := mock.NewMockReminderRepository(ctrl)
			mockTx := txMock.NewMockManager(ctrl)
			mockUUID := mockUUID.NewMockUUIDManager(ctrl)
			mockTime := mockTime.NewMockRealTimeProvider(ctrl)

			tt.prepareFunc(mockReminderRepo, mockUUID, mockTx, mockTime)
			uc := NewReminderUsecase(mockTx, mockUUID, mockReminderRepo, mockTime)

			ctx = context.WithValue(ctx, middleware.UserIDKey, testUserId)
			actualRes, err := uc.GetReminder(ctx, tt.reminderId)
			if tt.wantedError != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantedError)
				}
				assert.EqualError(t, err, tt.wantedError.Error())
				return
			}

			assert.Equal(t, tt.wantedRes, actualRes)
		})
	}
}

func TestGetReminders(t *testing.T) {
	logger.InitForTest()
	ctx := context.Background()
	testUserId := "test-user-id"
	testReminderId := "test-reminder-uuid"
	secondTestReminderId := "test-reminder-uuid"
	testTitle := "test title"
	testDescription := "test description"
	testLatitude := 20.00
	testLongitude := 20.00

	tests := []struct {
		name        string
		status      string
		prepareFunc func(
			mrr *mock.MockReminderRepository,
			mUUID *mockUUID.MockUUIDManager,
			mTx *txMock.MockManager,
			mTime *mockTime.MockRealTimeProvider,
		)
		wantedError error
		wantedRes   []domain.Reminder
	}{
		{
			name: "success when status is not provided",
			prepareFunc: func(
				mrr *mock.MockReminderRepository,
				mUUID *mockUUID.MockUUIDManager,
				mTx *txMock.MockManager,
				mTime *mockTime.MockRealTimeProvider,
			) {
				mrr.EXPECT().GetReminders(gomock.Any(), testUserId, "").Return(
					[]domain.Reminder{
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
						{
							ReminderId:  secondTestReminderId,
							UserId:      testUserId,
							Title:       testTitle,
							Description: testDescription,
							Latitude:    testLatitude,
							Longitude:   testLongitude,
							Radius:      domain.DefaultRadius,
							Status:      string(domain.PausedStatus),
						},
					}, nil,
				)
			},
			wantedError: nil,
			wantedRes: []domain.Reminder{
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
				{
					ReminderId:  secondTestReminderId,
					UserId:      testUserId,
					Title:       testTitle,
					Description: testDescription,
					Latitude:    testLatitude,
					Longitude:   testLongitude,
					Radius:      domain.DefaultRadius,
					Status:      string(domain.PausedStatus),
				},
			},
		},
		{
			name:   "fail when invalid status is provided",
			status: "invalid-status",
			prepareFunc: func(
				mrr *mock.MockReminderRepository,
				mUUID *mockUUID.MockUUIDManager,
				mTx *txMock.MockManager,
				mTime *mockTime.MockRealTimeProvider,
			) {
			},
			wantedError: apperror.BadRequest(),
			wantedRes:   []domain.Reminder{},
		},
		{
			name:   "success when status is completed",
			status: string(domain.CompletedStatus),
			prepareFunc: func(
				mrr *mock.MockReminderRepository,
				mUUID *mockUUID.MockUUIDManager,
				mTx *txMock.MockManager,
				mTime *mockTime.MockRealTimeProvider,
			) {
				mrr.EXPECT().GetReminders(gomock.Any(), testUserId, string(domain.CompletedStatus)).Return(
					[]domain.Reminder{
						{
							ReminderId:  testReminderId,
							UserId:      testUserId,
							Title:       testTitle,
							Description: testDescription,
							Latitude:    testLatitude,
							Longitude:   testLongitude,
							Radius:      domain.DefaultRadius,
							Status:      string(domain.CompletedStatus),
						},
						{
							ReminderId:  secondTestReminderId,
							UserId:      testUserId,
							Title:       testTitle,
							Description: testDescription,
							Latitude:    testLatitude,
							Longitude:   testLongitude,
							Radius:      domain.DefaultRadius,
							Status:      string(domain.CompletedStatus),
						},
					}, nil,
				)
			},
			wantedError: nil,
			wantedRes: []domain.Reminder{
				{
					ReminderId:  testReminderId,
					UserId:      testUserId,
					Title:       testTitle,
					Description: testDescription,
					Latitude:    testLatitude,
					Longitude:   testLongitude,
					Radius:      domain.DefaultRadius,
					Status:      string(domain.CompletedStatus),
				},
				{
					ReminderId:  secondTestReminderId,
					UserId:      testUserId,
					Title:       testTitle,
					Description: testDescription,
					Latitude:    testLatitude,
					Longitude:   testLongitude,
					Radius:      domain.DefaultRadius,
					Status:      string(domain.CompletedStatus),
				},
			},
		},
		{
			name:   "success when status is paused",
			status: string(domain.PausedStatus),
			prepareFunc: func(
				mrr *mock.MockReminderRepository,
				mUUID *mockUUID.MockUUIDManager,
				mTx *txMock.MockManager,
				mTime *mockTime.MockRealTimeProvider,
			) {
				mrr.EXPECT().GetReminders(gomock.Any(), testUserId, string(domain.PausedStatus)).Return(
					[]domain.Reminder{
						{
							ReminderId:  testReminderId,
							UserId:      testUserId,
							Title:       testTitle,
							Description: testDescription,
							Latitude:    testLatitude,
							Longitude:   testLongitude,
							Radius:      domain.DefaultRadius,
							Status:      string(domain.PausedStatus),
						},
						{
							ReminderId:  secondTestReminderId,
							UserId:      testUserId,
							Title:       testTitle,
							Description: testDescription,
							Latitude:    testLatitude,
							Longitude:   testLongitude,
							Radius:      domain.DefaultRadius,
							Status:      string(domain.PausedStatus),
						},
					}, nil,
				)
			},
			wantedError: nil,
			wantedRes: []domain.Reminder{
				{
					ReminderId:  testReminderId,
					UserId:      testUserId,
					Title:       testTitle,
					Description: testDescription,
					Latitude:    testLatitude,
					Longitude:   testLongitude,
					Radius:      domain.DefaultRadius,
					Status:      string(domain.PausedStatus),
				},
				{
					ReminderId:  secondTestReminderId,
					UserId:      testUserId,
					Title:       testTitle,
					Description: testDescription,
					Latitude:    testLatitude,
					Longitude:   testLongitude,
					Radius:      domain.DefaultRadius,
					Status:      string(domain.PausedStatus),
				},
			},
		},
		{
			name:   "fail when invalid status is provided",
			status: string(domain.ActiveStatus),
			prepareFunc: func(
				mrr *mock.MockReminderRepository,
				mUUID *mockUUID.MockUUIDManager,
				mTx *txMock.MockManager,
				mTime *mockTime.MockRealTimeProvider,
			) {
				mrr.EXPECT().GetReminders(gomock.Any(), testUserId, string(domain.ActiveStatus)).Return([]domain.Reminder{}, apperror.Internal())
			},
			wantedError: apperror.Internal(),
			wantedRes:   []domain.Reminder{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockReminderRepo := mock.NewMockReminderRepository(ctrl)
			mockTx := txMock.NewMockManager(ctrl)
			mockUUID := mockUUID.NewMockUUIDManager(ctrl)
			mockTime := mockTime.NewMockRealTimeProvider(ctrl)

			tt.prepareFunc(mockReminderRepo, mockUUID, mockTx, mockTime)
			uc := NewReminderUsecase(mockTx, mockUUID, mockReminderRepo, mockTime)

			ctx = context.WithValue(ctx, middleware.UserIDKey, testUserId)
			actualRes, err := uc.GetReminders(ctx, tt.status)
			if tt.wantedError != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantedError)
				}
				assert.EqualError(t, err, tt.wantedError.Error())
				return
			}

			assert.Equal(t, tt.wantedRes, actualRes)
		})
	}
}
