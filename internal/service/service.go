package service

import (
	"log"

	"github.com/theborzet/time-tracker/internal/config"
	"github.com/theborzet/time-tracker/internal/repository"
	externalapi "github.com/theborzet/time-tracker/pkg/external_api"
)

type ApiService struct {
	repo   *repository.ApiRepository
	logger *log.Logger
	cfg    *config.Config
	exApi  *externalapi.ExternalApiClient
}

func NewApiService(repo *repository.ApiRepository, logger *log.Logger, cfg *config.Config) *ApiService {
	client := externalapi.NewExternalApiClient(cfg.ExternalApi.APIURL)
	return &ApiService{
		repo:   repo,
		logger: logger,
		cfg:    cfg,
		exApi:  client,
	}
}
