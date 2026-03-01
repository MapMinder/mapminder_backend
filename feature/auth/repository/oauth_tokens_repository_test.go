package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MapMinder/mapminder_backend/feature/auth/domain"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	logger.InitForTest()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}

	dialector := mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm DB: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return gormDB, mock, cleanup
}

func TestOauthTokenRepository_GetOauthInformation(t *testing.T) {
	tests := []struct {
		name       string
		inputSub   string
		inputToken domain.OauthToken
		wantToken  domain.OauthToken
		wantErr    bool
		setupMock  func(mock sqlmock.Sqlmock)
	}{
		{
			name:     "GetOauthInformation returns existing record",
			inputSub: "google-subject",
			wantToken: domain.OauthToken{
				OauthId: 1, UserId: "user-123", OauthProvider: "google", OauthProviderId: "google-subject",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `oauth_token` WHERE oauth_provider_id = ? LIMIT ?")).
					WithArgs("google-subject", 1).
					WillReturnRows(sqlmock.NewRows([]string{"oauth_id", "user_id", "oauth_provider", "oauth_provider_id"}).
						AddRow(1, "user-123", "google", "google-subject"))
			},
		},
		{
			name:      "GetOauthInformation returns no record",
			inputSub:  "non-existent-sub",
			wantToken: domain.OauthToken{},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `oauth_token` WHERE oauth_provider_id = ? LIMIT ?")).
					WithArgs("non-existent-sub", 1).
					WillReturnRows(sqlmock.NewRows([]string{"oauth_id", "user_id", "oauth_provider", "oauth_provider_id"}))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()

			tt.setupMock(mock)

			repo := NewOauthTokenRepository(db)

			var got domain.OauthToken
			var err error
			got, err = repo.GetOauthInformation(context.Background(), tt.inputSub)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.inputSub != "" {
				assert.Equal(t, tt.wantToken, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOauthRepository_CreateOauthToken(t *testing.T) {
	tests := []struct {
		name       string
		inputSub   string
		inputToken domain.OauthToken
		wantToken  domain.OauthToken
		wantErr    bool
		setupMock  func(mock sqlmock.Sqlmock)
	}{
		{
			name: "CreateOauthToken succeeds",
			inputToken: domain.OauthToken{
				OauthId: 1, UserId: "user-456", OauthProvider: "google", OauthProviderId: "new-subject",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(
					"INSERT INTO `oauth_token` (`user_id`,`oauth_provider`,`oauth_provider_id`,`oauth_id`) VALUES (?,?,?,?)")).
					WithArgs("user-456", "google", "new-subject", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()

			tt.setupMock(mock)

			repo := NewOauthTokenRepository(db)

			var got domain.OauthToken
			var err error
			err = repo.CreateOauthToken(context.Background(), tt.inputToken)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.inputSub != "" {
				assert.Equal(t, tt.wantToken, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
