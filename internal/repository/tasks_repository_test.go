package repository

import (
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/theborzet/time-tracker/internal/models"
)

func TestGetUserTasks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	logger := log.New(os.Stdout, "", log.LstdFlags)
	repo := NewApiRepository(db, logger)

	userID := 1
	start := "2023-01-01"
	end := "2023-12-31"

	expectedTasks := []*models.TaskTimeSpent{
		{TaskName: "task1", TimeSpent: 60},
		{TaskName: "task2", TimeSpent: 30},
	}

	rows := sqlmock.NewRows([]string{"taskName", "timeSpent"}).
		AddRow("task1", 60).
		AddRow("task2", 30)

	mock.ExpectQuery(`SELECT taskName, EXTRACT\(EPOCH FROM \(endTime - startTime\)\)\/60 as timeSpent FROM tasks WHERE userId = \$1 AND startTime >= \$2 AND endTime <= \$3 ORDER BY timeSpent DESC`).
		WithArgs(userID, start, end).
		WillReturnRows(rows)

	tasks, err := repo.GetUserTasks(userID, start, end)
	assert.NoError(t, err)
	assert.NotNil(t, tasks)
	assert.Equal(t, len(expectedTasks), len(tasks))

	for i := range tasks {
		assert.Equal(t, expectedTasks[i].TaskName, tasks[i].TaskName)
		assert.Equal(t, expectedTasks[i].TimeSpent, tasks[i].TimeSpent)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestStartTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a stub database connection", err)
	}
	defer db.Close()

	logger := log.New(os.Stdout, "", log.LstdFlags)
	repo := NewApiRepository(db, logger)

	userID := 1
	taskName := "Task 1"
	startTime := "2023-07-15 09:00:00"

	mock.ExpectExec(`INSERT INTO tasks (.+) VALUES (.+)$`).
		WithArgs(userID, taskName, startTime).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.StartTask(userID, taskName, startTime)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEndTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a stub database connection", err)
	}
	defer db.Close()

	logger := log.New(os.Stdout, "", log.LstdFlags)
	repo := NewApiRepository(db, logger)

	userId := 1
	taskName := "Полить цветы"
	endTime := "2023-07-15 10:00:00"

	mock.ExpectExec(`UPDATE tasks SET endTime = (.+) WHERE userId = (.+) AND taskName = (.+) AND endTime IS NULL`).
		WithArgs(endTime, userId, taskName).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.EndTask(userId, taskName, endTime)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
