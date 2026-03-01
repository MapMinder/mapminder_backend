package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MapMinder/mapminder_backend/feature/reminder/domain"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	apperror "github.com/MapMinder/mapminder_backend/shared/appError"
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

func TestReminderRepository_Create(t *testing.T) {
	logger.InitForTest()
	testUserId := "test-user-id"
	testReminderId := "test-reminder-uuid"
	testTitle := "test title"
	testDescription := "test description"
	testLatitude := 20.00
	testLongitude := 20.00

	args := domain.Reminder{
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
		name      string
		args      domain.Reminder
		wantErr   bool
		setupMock func(mock sqlmock.Sqlmock)
	}{
		{
			name: "create reminders success",
			args: args,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `reminder` (`reminder_id`,`user_id`,`title`,`description`,`latitude`,`longitude`,`radius`,`status`,`last_triggered_at`,`completed_at`) VALUES (?,?,?,?,?,?,?,?,?,?)")).
					WithArgs(testReminderId, testUserId, testTitle, testDescription, testLatitude, testLongitude, domain.DefaultRadius, domain.CreatedStatus, nil, nil).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "returns error",
			args: args,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `reminder` (`reminder_id`,`user_id`,`title`,`description`,`latitude`,`longitude`,`radius`,`status`,`last_triggered_at`,`completed_at`) VALUES (?,?,?,?,?,?,?,?,?,?)")).
					WithArgs(testReminderId, testUserId, testTitle, testDescription, testLatitude, testLongitude, domain.DefaultRadius, domain.CreatedStatus, nil, nil).
					WillReturnError(apperror.Internal())
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()

			tt.setupMock(mock)

			repo := NewReminderRepository(db)
			err := repo.Create(context.Background(), tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
