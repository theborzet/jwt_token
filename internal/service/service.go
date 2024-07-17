package service

import (
	"log"

	"github.com/theborzet/time-tracker/config"
	"github.com/theborzet/time-tracker/internal/repository"
	externalapi "github.com/theborzet/time-tracker/pkg/external_api"
)

const (
	DateFormat      = "2006-01-02 15:04:05"
	TimestampFormat = "2006-01-02 15:04:05"
)

type ApiService struct {
	repo   *repository.ApiRepository
	logger *log.Logger
	cfg    *config.Config
	exApi  *externalapi.ExternalApiClient
}

func NewApiService(repo *repository.ApiRepository, logger *log.Logger, cfg *config.Config) *ApiService {
	client := externalapi.NewExternalApiClient(cfg.ExternalApiURL)
	return &ApiService{
		repo:   repo,
		logger: logger,
		cfg:    cfg,
		exApi:  client,
	}
}
