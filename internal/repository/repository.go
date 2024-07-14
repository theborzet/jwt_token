package repository

import (
	"github.com/theborzet/time-tracker/internal/models"

	"github.com/jmoiron/sqlx"
)

const (
	DateFormat      = "02.01.06 15:04:05"
	TimestampFormat = "2006-01-02 15:04:05"
)

type Repository interface {
	GetUsers(filter map[string]string) ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
	GetUserTasks(userID int, start, end string) ([]models.Task, error)
	StartTask(userID int, taskName string) error
	EndTask(userID int) error
}

type ApiRepository struct {
	db *sqlx.DB
}

func NewApiRepository(db *sqlx.DB) *ApiRepository {
	return &ApiRepository{
		db: db,
	}
}
