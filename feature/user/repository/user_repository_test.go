package repository_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MapMinder/mapminder_backend/feature/user/domain"
	"github.com/MapMinder/mapminder_backend/feature/user/repository"
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

func TestUserRepository_Create(t *testing.T) {
	tests := []struct {
		name      string
		user      domain.User
		wantErr   bool
		setupMock func(mock sqlmock.Sqlmock)
	}{
		{
			name: "Create user succeeds",
			user: domain.User{
				UserId:   "user-1234",
				Username: "test username",
				Email:    "test@gmail.com",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(
					"INSERT INTO `user` (`user_id`,`username`,`email`) VALUES (?,?,?)")).
					WithArgs("user-1234", "test username", "test@gmail.com").
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

			repo := repository.NewUserRepository(db)

			var err error
			err = repo.Create(context.Background(), tt.user)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetUser(t *testing.T) {
	tests := []struct {
		name      string
		userId    string
		wantUser  domain.User
		wantErr   bool
		setupMock func(mock sqlmock.Sqlmock)
	}{
		{
			name:     "Getuser returns existing record",
			userId:   "test-user-id",
			wantUser: domain.User{UserId: "test-user-id", Email: "test@gmail.com", Username: "test-username"},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `user` WHERE user_id = ? LIMIT ?")).
					WithArgs("test-user-id", 1).
					WillReturnRows(sqlmock.NewRows([]string{"test-user-id", "test@gmail.com", "test-username"}).
						AddRow("test-user-id", "test@gmail.com", "test-username"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()

			tt.setupMock(mock)

			repo := repository.NewUserRepository(db)

			var err error
			_, err = repo.GetUser(tt.userId)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
