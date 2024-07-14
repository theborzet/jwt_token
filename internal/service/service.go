package service

import "github.com/theborzet/time-tracker/internal/repository"

type ApiService struct {
	repo *repository.ApiRepository
}

func NewApiService(repo *repository.ApiRepository) *ApiService {
	return &ApiService{
		repo: repo,
	}
}
