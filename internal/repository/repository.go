package repository

import (
	"database/sql"
	"log"
	"time"
)

type Repository interface {
	SaveRefreshTokenHash(userID, hash string, expiresAt time.Time) error
	GetRefreshTokenHash(userID string) (string, time.Time, error)
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
