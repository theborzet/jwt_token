package service

import (
	"log"

	"github.com/theborzet/time-tracker/internal/repository"
)

type ApiService struct {
	repo   *repository.ApiRepository
	logger *log.Logger
}

func NewApiService(repo *repository.ApiRepository, logger *log.Logger) *ApiService {
	return &ApiService{
		repo:   repo,
		logger: logger,
	}
}
