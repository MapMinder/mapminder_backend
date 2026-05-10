package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

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
	testLocationName := "test location name"
	testLatitude := 20.00
	testLongitude := 20.00

	args := domain.Reminder{
		ReminderId:   testReminderId,
		UserId:       testUserId,
		Title:        testTitle,
		Description:  testDescription,
		Latitude:     testLatitude,
		Longitude:    testLongitude,
		LocationName: testLocationName,
		Radius:       domain.DefaultRadius,
		Status:       string(domain.ActiveStatus),
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
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `reminder` (`reminder_id`,`user_id`,`title`,`description`,`latitude`,`longitude`,`location_name`,`radius`,`status`,`last_triggered_at`,`completed_at`,`created_at`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")).
					WithArgs(testReminderId, testUserId, testTitle, testDescription, testLatitude, testLongitude, testLocationName, domain.DefaultRadius, domain.ActiveStatus, nil, nil, sqlmock.AnyArg()).
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
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `reminder` (`reminder_id`,`user_id`,`title`,`description`,`latitude`,`longitude`,`location_name`,`radius`,`status`,`last_triggered_at`,`completed_at`,`created_at`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")).
					WithArgs(testReminderId, testUserId, testTitle, testDescription, testLatitude, testLongitude, testLocationName, domain.DefaultRadius, domain.ActiveStatus, nil, nil, sqlmock.AnyArg()).
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

func TestReminderRepository_GetReminder(t *testing.T) {
	logger.InitForTest()
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
		name       string
		reminderId string
		wantErr    bool
		setupMock  func(mock sqlmock.Sqlmock)
	}{
		{
			name:       "get reminder success",
			reminderId: testReminderId,
			wantErr:    false,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"reminder_id", "user_id", "title", "description", "latitude", "longitude", "radius", "status", "last_triggered_at", "completed_at"}).
					AddRow(testReminderId, testUserId, testTitle, testDescription, testLatitude, testLongitude, domain.DefaultRadius, domain.ActiveStatus, nil, nil)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `reminder` WHERE reminder_id = ? LIMIT ?")).
					WithArgs(testReminderId, 1).
					WillReturnRows(rows)
			},
		},
		{
			name:       "get reminder not found",
			reminderId: testReminderId,
			wantErr:    true,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `reminder` WHERE reminder_id = ? LIMIT ?")).
					WithArgs(testReminderId, 1).
					WillReturnError(apperror.NotFound())
			},
		},
		{
			name:       "get reminder internal error",
			reminderId: testReminderId,
			wantErr:    true,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `reminder` WHERE reminder_id = ? LIMIT ?")).
					WithArgs(testReminderId, 1).
					WillReturnError(apperror.Internal())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()

			tt.setupMock(mock)

			repo := NewReminderRepository(db)
			reminder, err := repo.GetReminder(context.Background(), tt.reminderId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.name == "get reminder success" {
					assert.Equal(t, expectedReminder.ReminderId, reminder.ReminderId)
					assert.Equal(t, expectedReminder.UserId, reminder.UserId)
					assert.Equal(t, expectedReminder.Title, reminder.Title)
					assert.Equal(t, expectedReminder.Description, reminder.Description)
					assert.Equal(t, expectedReminder.Latitude, reminder.Latitude)
					assert.Equal(t, expectedReminder.Longitude, reminder.Longitude)
					assert.Equal(t, expectedReminder.Radius, reminder.Radius)
					assert.Equal(t, expectedReminder.Status, reminder.Status)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestReminderRepository_GetReminders(t *testing.T) {
	logger.InitForTest()
	testUserId := "test-user-id"
	secondTestReminderId := "test-remidner-uuid-2"
	testReminderId := "test-reminder-uuid"
	testTitle := "test title"
	testDescription := "test description"
	testLatitude := 20.00
	testLongitude := 20.00

	tests := []struct {
		name         string
		status       string
		userId       string
		wantErr      bool
		wantReminder []domain.Reminder
		setupMock    func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "success when no status is provided",
			userId:  testUserId,
			wantErr: false,
			wantReminder: []domain.Reminder{
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
					Status:      string(domain.ActiveStatus),
				},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"reminder_id", "user_id", "title", "description", "latitude", "longitude", "radius", "status", "last_triggered_at", "completed_at"}).
					AddRow(testReminderId, testUserId, testTitle, testDescription, testLatitude, testLongitude, domain.DefaultRadius, domain.ActiveStatus, nil, nil).
					AddRow(secondTestReminderId, testUserId, testTitle, testDescription, testLatitude, testLongitude, domain.DefaultRadius, domain.ActiveStatus, nil, nil)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `reminder` WHERE user_id = ?")).
					WithArgs(testUserId).
					WillReturnRows(rows)
			},
		},
		{
			name:    "when status active is provided",
			userId:  testUserId,
			status:  string(domain.ActiveStatus),
			wantErr: false,
			wantReminder: []domain.Reminder{
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
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"reminder_id", "user_id", "title", "description", "latitude", "longitude", "radius", "status", "last_triggered_at", "completed_at"}).
					AddRow(testReminderId, testUserId, testTitle, testDescription, testLatitude, testLongitude, domain.DefaultRadius, domain.ActiveStatus, nil, nil)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `reminder` WHERE user_id = ? AND status = ?")).
					WithArgs(testUserId, domain.ActiveStatus).
					WillReturnRows(rows)
			},
		},
		{
			name:    "when status paused is provided",
			userId:  testUserId,
			status:  string(domain.PausedStatus),
			wantErr: false,
			wantReminder: []domain.Reminder{
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
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"reminder_id", "user_id", "title", "description", "latitude", "longitude", "radius", "status", "last_triggered_at", "completed_at"}).
					AddRow(testReminderId, testUserId, testTitle, testDescription, testLatitude, testLongitude, domain.DefaultRadius, domain.PausedStatus, nil, nil)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `reminder` WHERE user_id = ? AND status = ?")).
					WithArgs(testUserId, domain.PausedStatus).
					WillReturnRows(rows)
			},
		},
		{
			name:    "when status completed is provided",
			userId:  testUserId,
			status:  string(domain.CompletedStatus),
			wantErr: false,
			wantReminder: []domain.Reminder{
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
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"reminder_id", "user_id", "title", "description", "latitude", "longitude", "radius", "status", "last_triggered_at", "completed_at"}).
					AddRow(testReminderId, testUserId, testTitle, testDescription, testLatitude, testLongitude, domain.DefaultRadius, domain.CompletedStatus, nil, nil)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `reminder` WHERE user_id = ? AND status = ?")).
					WithArgs(testUserId, domain.CompletedStatus).
					WillReturnRows(rows)
			},
		},
		{
			name:         "when status active is provided but no records match",
			userId:       testUserId,
			status:       string(domain.ActiveStatus),
			wantErr:      true,
			wantReminder: []domain.Reminder{},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `reminder` WHERE user_id = ? AND status = ?")).
					WithArgs(testUserId, domain.ActiveStatus).
					WillReturnError(apperror.NotFound())
			},
		},
		{
			name:         "internal error",
			userId:       testUserId,
			status:       string(domain.ActiveStatus),
			wantErr:      true,
			wantReminder: []domain.Reminder{},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `reminder` WHERE user_id = ? AND status = ?")).
					WithArgs(testUserId, domain.ActiveStatus).
					WillReturnError(apperror.Internal())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()

			tt.setupMock(mock)

			repo := NewReminderRepository(db)
			reminder, err := repo.GetReminders(context.Background(), tt.userId, tt.status)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantReminder, reminder)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestReminderRepository_DeleteReminder(t *testing.T) {
	logger.InitForTest()
	testReminderId := "test-reminder-id"
	tests := []struct {
		name       string
		reminderId string
		wantErr    bool
		setupMock  func(mock sqlmock.Sqlmock)
	}{
		{
			name:       "success",
			reminderId: testReminderId,
			wantErr:    false,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `reminder` WHERE reminder_id = ?").
					WithArgs(testReminderId).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:       "returns error when internal error occurs",
			reminderId: testReminderId,
			wantErr:    true,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `reminder` WHERE reminder_id = ?").
					WithArgs(testReminderId).
					WillReturnError(apperror.Internal())
				mock.ExpectRollback()
			},
		},
		{
			name:       "returns not found error when no rows affected",
			reminderId: testReminderId,
			wantErr:    true,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `reminder` WHERE reminder_id = ?").
					WithArgs(testReminderId).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()

			tt.setupMock(mock)

			repo := NewReminderRepository(db)
			err := repo.DeleteReminder(context.Background(), tt.reminderId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestReminderRepository_UpdateReminder(t *testing.T) {
	logger.InitForTest()
	testReminderId := "test-reminder-uuid"
	testTitle := "test title"
	testDescription := "test description"
	testLatitude := 20.00
	testLongitude := 20.00

	args := domain.Reminder{
		ReminderId:  testReminderId,
		Title:       testTitle,
		Description: testDescription,
		Latitude:    testLatitude,
		Longitude:   testLongitude,
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
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `reminder` SET `reminder_id`=?,`title`=?,`description`=?,`latitude`=?,`longitude`=? WHERE reminder_id = ?")).
					WithArgs(testReminderId, testTitle, testDescription, testLatitude, testLongitude, testReminderId).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "internal error",
			args: args,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `reminder` SET `reminder_id`=?,`title`=?,`description`=?,`latitude`=?,`longitude`=? WHERE reminder_id = ?")).
					WillReturnError(apperror.Internal())
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
			err := repo.UpdateReminder(context.Background(), tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestReminderRepository_UpdateLastTriggeredAt(t *testing.T) {
	logger.InitForTest()
	testReminderId := "test-reminder-uuid"
	now := time.Now()

	tests := []struct {
		name            string
		reminderId      string
		LastTriggeredAt time.Time
		wantErr         bool
		setupMock       func(mock sqlmock.Sqlmock)
	}{
		{
			name:            "create reminders success",
			reminderId:      testReminderId,
			LastTriggeredAt: now,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `reminder` SET `last_triggered_at`=? WHERE reminder_id = ?")).
					WithArgs(now, testReminderId).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:            "internal error",
			reminderId:      testReminderId,
			LastTriggeredAt: now,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `reminder` SET `last_triggered_at`=? WHERE reminder_id = ?")).
					WillReturnError(apperror.Internal())
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
			err := repo.UpdateLastTriggeredAt(context.Background(), tt.reminderId, tt.LastTriggeredAt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
