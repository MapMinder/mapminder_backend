package usecase

import (
	"context"
	"testing"

	"github.com/MapMinder/mapminder_backend/feature/auth/domain"
	"github.com/MapMinder/mapminder_backend/feature/auth/usecase/mock"
	usrDomain "github.com/MapMinder/mapminder_backend/feature/user/domain"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	apperror "github.com/MapMinder/mapminder_backend/shared/appError"
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

func TestGoogleSignIn(t *testing.T) {
	logger.InitForTest()
	t.Setenv("GOOGLE_CLIENT_ID", "test-client-id")
	t.Setenv("JWT_SIGN_KEY", "test-secret")
	ctx := context.Background()
	testSubject := "test-subject"
	testIdToken := "test-id-token"
	testClientId := "test-client-id"
	testCreds := domain.GoogleSignInCreds{
		IdToken:  testIdToken,
		ClientId: testClientId,
	}
	testUserName := "test username"
	testEmail := "testemail@gmail.com"
	testClaims := domain.GoogleIdToken{
		Subject: testSubject,
		Claims: domain.GoogleIdTokenClaims{
			Name:  testUserName,
			Email: testEmail,
		},
	}

	testUserUUID := "test-user-uuid"

	testUser := usrDomain.User{
		UserId:   testUserUUID,
		Email:    testClaims.Claims.Email,
		Username: testClaims.Claims.Name,
	}

	testOauthToken := domain.OauthToken{
		UserId:          testUserUUID,
		OauthProvider:   domain.GoogleProvider,
		OauthProviderId: testClaims.Subject,
	}

	tests := []struct {
		name        string
		prepareFunc func(
			mgr *mock.MockGoogleSignInApiRepository,
			motr *mock.MockOauthTokenRepository,
			mur *mock.MockUserRepository,
			mUUID *mockUUID.MockUUIDManager,
			mTx *txMock.MockManager,
		)
		wantedError error
	}{
		{
			name: "Success when user already exists",
			prepareFunc: func(
				mgr *mock.MockGoogleSignInApiRepository,
				motr *mock.MockOauthTokenRepository,
				mur *mock.MockUserRepository,
				mUUID *mockUUID.MockUUIDManager,
				mTx *txMock.MockManager,
			) {
				mgr.EXPECT().GetGoogleUserInfo(gomock.Any(), testCreds).Return(testClaims, nil)
				motr.EXPECT().GetOauthInformation(gomock.Any(), testClaims.Claims.Subject).Return(testOauthToken, nil)
				mur.EXPECT().GetUser(gomock.Any(), testUserUUID).Return(testUser, nil)
			},
			wantedError: nil,
		},
		{
			name: "Success when user does not exist",
			prepareFunc: func(
				mgr *mock.MockGoogleSignInApiRepository,
				motr *mock.MockOauthTokenRepository,
				mur *mock.MockUserRepository,
				mUUID *mockUUID.MockUUIDManager,
				mTx *txMock.MockManager,
			) {
				mgr.EXPECT().GetGoogleUserInfo(gomock.Any(), testCreds).Return(testClaims, nil)
				motr.EXPECT().GetOauthInformation(gomock.Any(), testClaims.Claims.Subject).Return(domain.OauthToken{}, nil)
				mUUID.EXPECT().NewV7().Return(testUserUUID, nil)
				mur.EXPECT().Create(gomock.Any(), testUser).Return(nil)
				motr.EXPECT().CreateOauthToken(gomock.Any(), testOauthToken).Return(nil)
			},
			wantedError: nil,
		},
		{
			name: "Fails to get google user information",
			prepareFunc: func(
				mgr *mock.MockGoogleSignInApiRepository,
				motr *mock.MockOauthTokenRepository,
				mur *mock.MockUserRepository,
				mUUID *mockUUID.MockUUIDManager,
				mTx *txMock.MockManager,
			) {
				mgr.EXPECT().GetGoogleUserInfo(gomock.Any(), testCreds).Return(domain.GoogleIdToken{}, apperror.GoogleError())
			},
			wantedError: apperror.GoogleError(),
		},
		{
			name: "Fails to marshal or unmarshal",
			prepareFunc: func(
				mgr *mock.MockGoogleSignInApiRepository,
				motr *mock.MockOauthTokenRepository,
				mur *mock.MockUserRepository,
				mUUID *mockUUID.MockUUIDManager,
				mTx *txMock.MockManager,
			) {
				mgr.EXPECT().GetGoogleUserInfo(gomock.Any(), testCreds).Return(domain.GoogleIdToken{}, apperror.Internal())
			},
			wantedError: apperror.Internal(),
		},
		{
			name: "Fails to get OauthInformation",
			prepareFunc: func(
				mgr *mock.MockGoogleSignInApiRepository,
				motr *mock.MockOauthTokenRepository,
				mur *mock.MockUserRepository,
				mUUID *mockUUID.MockUUIDManager,
				mTx *txMock.MockManager,
			) {
				mgr.EXPECT().GetGoogleUserInfo(gomock.Any(), testCreds).Return(testClaims, nil)
				motr.EXPECT().GetOauthInformation(gomock.Any(), testClaims.Claims.Subject).Return(domain.OauthToken{}, apperror.Internal())
			},
			wantedError: apperror.Internal(),
		},
		{
			name: "Fails to create user",
			prepareFunc: func(
				mgr *mock.MockGoogleSignInApiRepository,
				motr *mock.MockOauthTokenRepository,
				mur *mock.MockUserRepository,
				mUUID *mockUUID.MockUUIDManager,
				mTx *txMock.MockManager,
			) {
				mgr.EXPECT().GetGoogleUserInfo(gomock.Any(), testCreds).Return(testClaims, nil)
				motr.EXPECT().GetOauthInformation(gomock.Any(), testClaims.Claims.Subject).Return(domain.OauthToken{}, nil)
				mUUID.EXPECT().NewV7().Return(testUserUUID, nil)
				mur.EXPECT().Create(gomock.Any(), testUser).Return(apperror.Internal())
			},
			wantedError: apperror.Internal(),
		},
		{
			name: "Fails to create oauth information",
			prepareFunc: func(
				mgr *mock.MockGoogleSignInApiRepository,
				motr *mock.MockOauthTokenRepository,
				mur *mock.MockUserRepository,
				mUUID *mockUUID.MockUUIDManager,
				mTx *txMock.MockManager,
			) {
				mgr.EXPECT().GetGoogleUserInfo(gomock.Any(), testCreds).Return(testClaims, nil)
				motr.EXPECT().GetOauthInformation(gomock.Any(), testClaims.Claims.Subject).Return(domain.OauthToken{}, nil)
				mUUID.EXPECT().NewV7().Return(testUserUUID, nil)
				mur.EXPECT().Create(gomock.Any(), testUser).Return(nil)
				motr.EXPECT().CreateOauthToken(gomock.Any(), testOauthToken).Return(apperror.Internal())
			},
			wantedError: apperror.Internal(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockGoogle := mock.NewMockGoogleSignInApiRepository(ctrl)
			mockOauth := mock.NewMockOauthTokenRepository(ctrl)
			mockUser := mock.NewMockUserRepository(ctrl)
			mockTx := txMock.NewMockManager(ctrl)
			mockUUID := mockUUID.NewMockUUIDManager(ctrl)

			// setup
			mockWithinTransaction(t, mockTx)

			tt.prepareFunc(mockGoogle, mockOauth, mockUser, mockUUID, mockTx)
			uc := NewGoogleSignInUsecase(mockOauth, mockGoogle, mockUser, mockUUID, mockTx)

			token, err := uc.GoogleSignIn(ctx, "test-id-token")
			if tt.wantedError != nil {
				// test expects an error
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantedError)
				}
				assert.EqualError(t, err, tt.wantedError.Error())
				return
			}

			// test expects success
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			assert.NotEmpty(t, token)
		})
	}
}
