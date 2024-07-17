package repository

import (
	"database/sql"
	"log"

	"github.com/theborzet/time-tracker/internal/models"
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
	db     *sql.DB
	logger *log.Logger
}

func NewApiRepository(db *sql.DB, logger *log.Logger) *ApiRepository {
	return &ApiRepository{
		db:     db,
		logger: logger,
	}
}
