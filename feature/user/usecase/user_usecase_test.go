package usecase

import (
	"context"
	"testing"

	"github.com/MapMinder/mapminder_backend/feature/user/domain"
	"github.com/MapMinder/mapminder_backend/feature/user/usecase/mock"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	apperror "github.com/MapMinder/mapminder_backend/shared/appError"
	"github.com/MapMinder/mapminder_backend/shared/middleware"
	txMock "github.com/MapMinder/mapminder_backend/shared/tx/mock"
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

func TestUserUsecase_GetUser(t *testing.T) {
	logger.InitForTest()
	ctx := context.Background()
	testUserId := "test-user-id"

	testUser := domain.User{
		UserId:   testUserId,
		Username: "testuser",
		Email:    "test@example.com",
	}

	tests := []struct {
		name        string
		prepareFunc func(
			mur *mock.MockUserRepository,
			mTx *txMock.MockManager,
		)
		wantedError error
		wantedRes   domain.User
	}{
		{
			name: "success",
			prepareFunc: func(
				mur *mock.MockUserRepository,
				mTx *txMock.MockManager,
			) {
				mur.EXPECT().GetUser(gomock.Any(), testUserId).Return(testUser, nil)
			},
			wantedError: nil,
			wantedRes:   testUser,
		},
		{
			name: "fails when repository returns error",
			prepareFunc: func(
				mur *mock.MockUserRepository,
				mTx *txMock.MockManager,
			) {
				mur.EXPECT().GetUser(gomock.Any(), testUserId).Return(domain.User{}, apperror.Internal())
			},
			wantedError: apperror.Internal(),
			wantedRes:   domain.User{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mock.NewMockUserRepository(ctrl)
			mockTx := txMock.NewMockManager(ctrl)

			mockWithinTransaction(t, mockTx)

			tt.prepareFunc(mockUserRepo, mockTx)
			uc := NewUserUsecase(mockTx, mockUserRepo)

			ctx = context.WithValue(ctx, middleware.UserIDKey, testUserId)
			actualRes, err := uc.GetUser(ctx)
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

